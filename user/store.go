package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	SubjectiveStyleField          = "subjectiveStyle"
	ObjectiveStyleField           = "objectiveStyle"
	AssessmentAndPlanStyleField   = "assessmentStyle"
	PlanningStyleField            = "planningStyle"
	SummaryStyleField             = "summaryStyle"
	PatientInstructionsStyleField = "patientInstructionsStyle"
)

var EmailAlreadyExistsError = errors.New("email already exists")

const EmailField = "email"

type User struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty"`
	Name                     string             `bson:"name"`
	Email                    string             `bson:"email"`
	PasswordHash             string             `bson:"passwordHash"`
	SubjectiveStyle          string             `bson:"subjectiveStyle"`
	ObjectiveStyle           string             `bson:"objectiveStyle"`
	AssessmentAndPlanStyle   string             `bson:"assessmentStyle"`
	SummaryStyle             string             `bson:"summaryStyle"`
	PatientInstructionsStyle string             `bson:"patientInstructionsStyle"`
}

type UserStore interface {
	Put(ctx context.Context, name, email, password string) (string, error)
	Get(ctx context.Context, id string) (User, error)
	GetByAuth(ctx context.Context, email, password string) (User, error)
	GetStyleField(ctx context.Context, userID, styleField string) (string, error)
	UpdateStyle(ctx context.Context, providerID, contentType, newStyle string) error
	UpdateProfileSettings(ctx context.Context, userID string, name string, currentPassword string, newPassword string) error 
	CheckEmailExistence(ctx context.Context, email string) (bool, error)

}

type store struct {
	client *mongo.Collection
}

func NewUserStore(client *mongo.Collection) UserStore {
	// Create a unique index on the email field
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	indexExists, err := indexExists(ctx, client.Indexes(), EmailField)
	if err != nil {
		panic(fmt.Errorf("error checking if index on email field already exists: %v", err))
	}
	if !indexExists {
		model := mongo.IndexModel{
			Keys:    bson.D{{Key: EmailField, Value: 1}},
			Options: options.Index().SetUnique(true),
		}
		_, err := client.Indexes().CreateOne(ctx, model)
		if err != nil {
			panic(fmt.Errorf("error creating unique index on email field: %v", err))
		}
	}
	return &store{client: client}
}

func indexExists(ctx context.Context, indexView mongo.IndexView, fieldName string) (bool, error) {
	cursor, err := indexView.List(ctx)
	if err != nil {
		return false, fmt.Errorf("error listing indexes: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var indexInfo bson.M
		if err := cursor.Decode(&indexInfo); err != nil {
			return false, fmt.Errorf("error decoding index info: %v", err)
		}
		if name, ok := indexInfo["name"].(string); ok && name == fieldName {
			return true, nil
		}
	}
	return false, nil
}

func (s *store) CheckEmailExistence(ctx context.Context, email string) (bool, error) {
	filter := bson.M{EmailField: email}
	var existingUser User
	err := s.client.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, fmt.Errorf("error occurred while checking for email existence : %v", err)
	}
	return false, EmailAlreadyExistsError
}

func (s *store) Put(ctx context.Context, name, email, password string) (string, error) {
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
	case SubjectiveStyleField, ObjectiveStyleField, AssessmentAndPlanStyleField, PlanningStyleField, SummaryStyleField, PatientInstructionsStyleField:
		return true
	default:
		return false
	}
}

func (s *store) GetStyleField(ctx context.Context, userID, styleField string) (string, error) {
	if !IsValidStyle(styleField) {
		return "", fmt.Errorf("invalid style: %s", styleField)
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objectID}
	projection := bson.M{styleField: 1, "_id": 0}
	opts := options.FindOne().SetProjection(projection)

	var result bson.M
	err = s.client.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("user not found")
		}
		return "", fmt.Errorf("failed to fetch user: %v", err)
	}

	stringValue, ok := result[styleField].(string)
	if ok {
		return stringValue, nil
	} else {
		return "", fmt.Errorf("error: field '%s' not found or not a string", stringValue)
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


func (s *store) UpdateProfileSettings(ctx context.Context, userID string, name string, currentPassword string, newPassword string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objectID}
	var existingUser User
	err = s.client.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to retrieve user: %v", err)
	}

	update := bson.M{}

	if name != existingUser.Name {
		update["name"] = name
	}

	if newPassword != "" {

		if currentPassword == newPassword{
			return errors.New("new password cannot be the same as current password")
		}

		err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(currentPassword))
		if err != nil {
			return errors.New("invalid current password")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing new password: %v", err)
		}
		update["passwordHash"] = string(hashedPassword)
	}

	if len(update) > 0 {
		_, err = s.client.UpdateOne(ctx, filter, bson.M{"$set": update})
		if err != nil {
			return fmt.Errorf("failed to update user profile: %v", err)
		}
	}

	return nil
}

