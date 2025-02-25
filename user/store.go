package user

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	SubjectiveStyleField = "subjectiveStyle"
	ObjectiveStyleField  = "objectiveStyle"
	AssessmentStyleField = "assessmentStyle"
	PlanningStyleField   = "planningStyle"
	SummaryStyleField    = "summaryStyle"
)

const EmailField = "email"

type Styles struct {
	SubjectiveStyle string `bson:"subjectiveStyle"`
	ObjectiveStyle  string `bson:"objectiveStyle"`
	AssessmentStyle string `bson:"assessmentStyle"`
	PlanningStyle   string `bson:"planningStyle"`
	SummaryStyle    string `bson:"summaryStyle"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"passwordHash"`
	Styles
}

type UserStore interface {
	Put(ctx context.Context, name, email, password string) (string, error)
	Get(ctx context.Context, id string) (User, error)
	GetByAuth(ctx context.Context, email, password string) (User, error)
	UpdateStyle(ctx context.Context, providerID, contentType, newStyle string) error
}

func NewUserStore(client *mongo.Collection) UserStore {
	return &store{client: client}
}

type store struct {
	client *mongo.Collection
}

func (s *store) Put(ctx context.Context, name, email, password string) (string, error) {
	filter := bson.M{EmailField: email}
	var existingUser User
	err := s.client.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return "", fmt.Errorf("user already exists with this email: %s", email)
	} else if err != mongo.ErrNoDocuments {
		return "", fmt.Errorf("failed to check for existing user: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password:%v", err)
	}

	newUser := User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	insertResp, err := s.client.InsertOne(ctx, newUser)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %v", err)
	}

	insertID, ok := insertResp.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("unexpected type for InsertedID: %T", insertID)
	}

	return insertID.Hex(), nil
}

func (s *store) Get(ctx context.Context, id string) (User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User{}, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objectID}
	var retrievedUser User
	err = s.client.FindOne(ctx, filter).Decode(&retrievedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, errors.New("user not found")
		}
		return User{}, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return retrievedUser, nil
}

func (s *store) GetByAuth(ctx context.Context, email, password string) (User, error) {

	var retrievedUser User
	filter := bson.M{EmailField: email}
	err := s.client.FindOne(ctx, filter).Decode(&retrievedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, errors.New("user not found")
		}
		return User{}, fmt.Errorf("failed to fetch user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(retrievedUser.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return User{}, fmt.Errorf("incorrect authentication credentials: %v", err)
		}
		return User{}, err
	}
	return retrievedUser, nil
}

func IsValidStyle(style string) bool {
	switch style {
	case SubjectiveStyleField, ObjectiveStyleField, AssessmentStyleField, PlanningStyleField, SummaryStyleField:
		return true
	default:
		return false
	}
}

func (s *store) UpdateStyle(ctx context.Context, providerID, styleField, newStyle string) error {

	objectID, err := primitive.ObjectIDFromHex(providerID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	if !IsValidStyle(styleField) {
		return fmt.Errorf("invalid style field: %s", styleField)
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: styleField, Value: newStyle}}}}

	result, err := s.client.UpdateOne(ctx, bson.D{{Key: "_id", Value: objectID}}, update)
	if err != nil {
		return fmt.Errorf("error updating the report field in MongoDB: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with id %s", providerID)
	}

	return nil
}
