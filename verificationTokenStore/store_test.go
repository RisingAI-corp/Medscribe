package verificationStore

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test Data
const (
	testToken           = "test_token_123"
	testEmail           = "test@example.com"
	testProviderID      = "provider_abc"
	testNewPasswordHash = "hashed_password"
	testName            = "Test User"
)

func loadEnv(t *testing.T) {
	t.Helper()
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func setupTestDB(t *testing.T) *mongo.Collection {
	t.Helper()
	loadEnv(t)

	mongoURI := os.Getenv("MONGODB_URI_DEV")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in .env file")
	}

	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		log.Fatal("MONGODB_DB is not set in .env file")
	}

	collectionName := os.Getenv("MONGODB_VERIFICATION_COLLECTION_DEV") // You might need to set this in your .env
	if collectionName == "" {
		log.Fatal("MONGODB_VERIFICATION_COLLECTION_DEV is not set in .env file")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return collection
}

func cleanupTestDB(t *testing.T, collection *mongo.Collection) {
	t.Helper()
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("Failed to clean up test database: %v", err)
	}
}

func TestGetBufferedUserDocument(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })
	store, err := NewVerificationStore(context.Background(), collection, 60) // Initialize the store correctly
	assert.NoError(t, err)
	ctx := context.Background()

	docToInsert := BufferedUserDocument{
		Token:     testToken,
		Email:     testEmail,
		CreatedAt: time.Now(),
	}

	_, err = collection.InsertOne(ctx, docToInsert)
	assert.NoError(t, err)

	t.Run("should return BufferedUserDocument when token exists", func(t *testing.T) {
		retrievedDoc, err := store.GetBufferedUserDocument(ctx, testToken)
		assert.NoError(t, err)
		assert.Equal(t, testToken, retrievedDoc.Token)
		assert.Equal(t, testEmail, retrievedDoc.Email)
	})

	t.Run("should return error when token does not exist", func(t *testing.T) {
		_, err := store.GetBufferedUserDocument(ctx, "non_existent_token")
		assert.Error(t, err)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
		assert.Contains(t, err.Error(), "verification token 'non_existent_token' not found")
	})
}

func TestGetResetPasswordDetails(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })
	store, err := NewVerificationStore(context.Background(), collection, 60) // Initialize the store correctly
	assert.NoError(t, err)
	ctx := context.Background()

	docToInsert := ResetPasswordDetails{
		Token:         testToken,
		Email:    testProviderID,
		NewPasswordHash: testNewPasswordHash,
		CreatedAt:     time.Now(),
	}

	_, err = collection.InsertOne(ctx, docToInsert)
	assert.NoError(t, err)

	t.Run("should return ResetPasswordDetails when token exists", func(t *testing.T) {
		retrievedDoc, err := store.GetResetPasswordDetails(ctx, testToken)
		assert.NoError(t, err)
		assert.Equal(t, testToken, retrievedDoc.Token)
		assert.Equal(t, testProviderID, retrievedDoc.Email)
		assert.Equal(t, testNewPasswordHash, retrievedDoc.NewPasswordHash)
	})

	t.Run("should return error when token does not exist", func(t *testing.T) {
		_, err := store.GetResetPasswordDetails(ctx, "non_existent_token")
		assert.Error(t, err)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
		assert.Contains(t, err.Error(), "verification token 'non_existent_token' not found")
	})
}

func TestPutBufferedUserDocument(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })
	store, err := NewVerificationStore(context.Background(), collection, 60) // Initialize the store correctly
	assert.NoError(t, err)
	ctx := context.Background()

	t.Run("should successfully insert a new BufferedUserDocument", func(t *testing.T) {
		err := store.PutBufferedUserDocument(ctx, testToken, testName, testEmail, testNewPasswordHash)
		assert.NoError(t, err)

		// Verify it was inserted
		var retrievedDoc BufferedUserDocument
		err = collection.FindOne(ctx, bson.M{"token": testToken}).Decode(&retrievedDoc)
		assert.NoError(t, err)
		assert.Equal(t, testToken, retrievedDoc.Token)
		assert.Equal(t, testEmail, retrievedDoc.Email)
		assert.False(t, retrievedDoc.CreatedAt.IsZero())
	})

	t.Run("should return error on duplicate token", func(t *testing.T) {
		err := store.PutBufferedUserDocument(ctx, testToken, testName, testEmail, testNewPasswordHash)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("duplicate token '%s'", testToken))
	})
}

func TestPutResetPasswordDetails(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })
	store, err := NewVerificationStore(context.Background(), collection, 60) // Initialize the store correctly
	assert.NoError(t, err)
	ctx := context.Background()

	t.Run("should successfully insert new ResetPasswordDetails", func(t *testing.T) {
		err := store.PutResetPasswordDetails(ctx, testToken, testProviderID)
		assert.NoError(t, err)

		// Verify it was inserted
		var retrievedDoc ResetPasswordDetails
		err = collection.FindOne(ctx, bson.M{"token": testToken}).Decode(&retrievedDoc)
		assert.NoError(t, err)
		assert.Equal(t, testToken, retrievedDoc.Token)
		assert.Equal(t, testProviderID, retrievedDoc.Email)
		assert.False(t, retrievedDoc.CreatedAt.IsZero())
	})

	t.Run("should return error on duplicate token", func(t *testing.T) {
		err := store.PutResetPasswordDetails(ctx, testToken, testProviderID) // Insert again
		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("duplicate token '%s'", testToken))
	})
}

func TestDelete(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })
	store, err := NewVerificationStore(context.Background(), collection, 60) // Initialize the store correctly
	assert.NoError(t, err)
	ctx := context.Background()

	docToInsert := BufferedUserDocument{
		Token: testToken,
		Email: testEmail,
	}
	_, err = collection.InsertOne(ctx, docToInsert)
	assert.NoError(t, err)

	t.Run("should successfully delete a document by token", func(t *testing.T) {
		err := store.Delete(ctx, testToken)
		assert.NoError(t, err)

		// Verify it's deleted
		count, err := collection.CountDocuments(ctx, bson.M{"token": testToken})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("should return error if token does not exist", func(t *testing.T) {
		err := store.Delete(ctx, "non_existent_token")
		assert.Error(t, err) // You're correct in this assertion
	})
}