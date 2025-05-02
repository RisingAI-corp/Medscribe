package verificationStore

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrVerificationDocumentNotFound is returned when the verification document is not found
var ErrVerificationDocumentNotFound = errors.New("verification document not found")
var ErrVerificationDocumentAlreadyExists = errors.New("verification document already exists")

// ForgotPasswordPayload holds data specific to password reset requests.
type ResetPasswordDetails struct {
	Token		  string `bson:"token"`          // The unique verification token (indexed).
	Email   string `bson:"email"`     // The provider ID associated with the token (indexed).
	NewPasswordHash string `bson:"newPasswordHash"` // The new hashed password to be set.
	CreatedAt time.Time        `bson:"createdAt"`  // Timestamp used for TTL expiration (TTL indexed).
}

// VerificationDocument represents the data stored for any verification code.
type BufferedUserDocument struct {
	Token     string           `bson:"token"`      // The unique verification token (indexed).
	Name 	string           `bson:"name"`       // The name associated with the token.
	Email     string           `bson:"email"`      // The email associated with the token (indexed). Used as the primary identifier.
	CreatedAt time.Time        `bson:"createdAt"`  // Timestamp used for TTL expiration (TTL indexed).
	Password   string      `bson:"password"`    // Holds either SignupPayload or ForgotPasswordPayload based on Type.
}


// VerificationStore defines the interface for verification code operations.
type VerificationStore interface {
    GetBufferedUserDocument(ctx context.Context, token string) (BufferedUserDocument, error)
    GetResetPasswordDetails(ctx context.Context, token string) (ResetPasswordDetails, error)
	PutBufferedUserDocument(ctx context.Context, token, name, email, password string) error
    PutResetPasswordDetails(ctx context.Context, token, email string) error
    Delete(ctx context.Context, token string) error
}

// verificationStore is the concrete implementation using MongoDB.
type verificationStore struct {
	coll *mongo.Collection
}

// NewVerificationStore creates a new VerificationStore instance and ensures necessary indexes exist.
func NewVerificationStore(ctx context.Context, collection *mongo.Collection, ttlSeconds int32) (VerificationStore, error) {
	// 1. Input Validation
	if err := validateNewVerificationStoreInput(collection, ttlSeconds); err != nil {
		return nil, err
	}

	// 2. Initialize Index Management
	indexView := collection.Indexes()

	// 3. Define a Helper Function for Checking Index Existence
	indexExists := func(indexName string) bool {
		return checkIndexExistence(ctx, indexView, indexName)
	}

	// 4. Ensure TTL Index on 'createdAt'
	if err := ensureTTLIndex(ctx, indexView, ttlSeconds, indexExists); err != nil {
		return nil, err
	}

	// 5. Ensure Unique Index on 'token'
	if err := ensureUniqueTokenIndex(ctx, indexView, indexExists); err != nil {
		return nil, fmt.Errorf("error ensuring unique token index: %w", err)
	}

	// 6. Ensure Unique Index on 'email'
	if err := ensureUniqueEmailIndex(ctx, indexView, indexExists); err != nil {
		return nil, fmt.Errorf("error ensuring unique email index: %w", err)
	}

	// 7. Return the Verification Store Instance
	return &verificationStore{coll: collection}, nil
}

// --- Helper Functions ---

// 1. Input Validation
func validateNewVerificationStoreInput(collection *mongo.Collection, ttlSeconds int32) error {
	if collection == nil {
		return errors.New("mongodb collection cannot be nil")
	}
	if ttlSeconds <= 0 {
		return errors.New("ttlSeconds must be positive")
	}
	return nil
}

// 2. Check Index Existence
func checkIndexExistence(ctx context.Context, indexView mongo.IndexView, indexName string) bool {
	cursor, err := indexView.List(ctx)
	if err != nil {
		log.Printf("Error listing indexes: %v", err)
		return false
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var indexInfo bson.M
		if err := cursor.Decode(&indexInfo); err != nil {
			log.Printf("Error decoding index info: %v", err)
			return false
		}
		if name, ok := indexInfo["name"].(string); ok && name == indexName {
			return true
		}
	}
	return false
}

// 4. Ensure TTL Index
func ensureTTLIndex(ctx context.Context, indexView mongo.IndexView, ttlSeconds int32, indexExists func(string) bool) error {
	ttlIndexName := "ttl_createdAt"
	if !indexExists(ttlIndexName) {
		ttlIndexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "createdAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(ttlSeconds).SetName(ttlIndexName),
		}
		_, err := indexView.CreateOne(ctx, ttlIndexModel)
		if err != nil {
			return fmt.Errorf("failed to create TTL index on createdAt: %w", err)
		}
	}
	return nil
}

// 5. Ensure Unique Token Index
func ensureUniqueTokenIndex(ctx context.Context, indexView mongo.IndexView, indexExists func(string) bool) error {
	tokenIndexName := "unique_token"
	if !indexExists(tokenIndexName) {
		tokenIndexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "token", Value: 1}},
			Options: options.Index().SetUnique(true).SetName(tokenIndexName),
		}
		_, err := indexView.CreateOne(ctx, tokenIndexModel)
		if err != nil {
			return fmt.Errorf("failed to create unique token index: %w", err)
		}
	}
	return nil
}

// 6. Ensure Unique Email Index
func ensureUniqueEmailIndex(ctx context.Context, indexView mongo.IndexView, indexExists func(string) bool) error {
	emailIndexName := "email_lookup"
	if !indexExists(emailIndexName) {
		emailIndexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName(emailIndexName),
		}
		_, err := indexView.CreateOne(ctx, emailIndexModel)
		if err != nil {
			return fmt.Errorf("failed to create unique email index: %w", err)
		}
	}
	return nil
}

// Get retrieves a verification document by its token. Caller must check Type and assert Payload type.
func (s *verificationStore) GetBufferedUserDocument(ctx context.Context, token string) (BufferedUserDocument, error) {
	var doc BufferedUserDocument
	filter := bson.M{"token": token}

	// Try with default registry first
	err := s.coll.FindOne(ctx, filter).Decode(&doc)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return BufferedUserDocument{}, ErrVerificationDocumentNotFound
	}
	if err != nil {
        return BufferedUserDocument{}, fmt.Errorf("failed to get verification document by token '%s': %w", token, err)
    }
	return doc, nil
}

func (s *verificationStore) GetResetPasswordDetails(ctx context.Context, token string) (ResetPasswordDetails, error) {
	var doc ResetPasswordDetails
	filter := bson.M{"token": token}

	err := s.coll.FindOne(ctx, filter).Decode(&doc)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return ResetPasswordDetails{}, fmt.Errorf("verification token '%s' not found: %w", token, err)
	}
	if err != nil {
        return ResetPasswordDetails{}, fmt.Errorf("failed to get verification document by token '%s': %w", token, err)
    }
	return doc, nil
}

func (s *verificationStore) PutBufferedUserDocument(ctx context.Context, token, name, email, password string) error {
	fmt.Println("Inserting verification document with token:", token)
	doc := BufferedUserDocument{
		Token:     token,
		Email:     email,
		Name: 	name,	
		Password:  password,
		CreatedAt: time.Now(),
	}

	_, err := s.coll.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrVerificationDocumentAlreadyExists
		}
		return fmt.Errorf("failed to insert verification document: %w", err)
	}
	return nil
}

// PutForgotPassword stores a new verification document for forgot password.
func (s *verificationStore) PutResetPasswordDetails(ctx context.Context, token, email string) error {
	doc := ResetPasswordDetails{
		Token:         token,
		Email:    email,
		CreatedAt:     time.Now(),
	}
	_, err := s.coll.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrVerificationDocumentAlreadyExists
		}
		return fmt.Errorf("failed to insert verification document: %w", err)
	}
	return nil
}

// Delete removes a verification document by its token.
func (s *verificationStore) Delete(ctx context.Context, token string) error {
	filter := bson.M{"token": token}
	result, err := s.coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete verification token '%s': %w", token, err)
	}
	if result.DeletedCount == 0 {
		return errors.New("no document found to delete")
	}
	return nil
}

func GenerateRandomToken(size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("number of digits must be positive: %d", size)
	}

	// Use a buffer for efficient string concatenation.
	var result string
	for i := 0; i < size; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))	
		if err != nil {
			return "", fmt.Errorf("failed to generate random digit: %w", err)
		}
		result += num.String()
	}
	fmt.Println("Generated token:", result)
	return result, nil
}
