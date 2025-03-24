package user

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const email = "john.doe@example.com"
const name = "John Doe"
const password = "password123"

func loadEnv(t *testing.T) {
	t.Helper()
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func setupTestDB(t *testing.T) *mongo.Collection {
	t.Helper()
	loadEnv(t)

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in .env file")
	}

	dbName := os.Getenv("MONGODB_DB")
	if mongoURI == "" {
		log.Fatal("MONGODB_DB is not set in .env file")
	}

	collectionName := os.Getenv("MONGODB_USER_COLLECTION_DEV")
	if mongoURI == "" {
		log.Fatal("MONGODB_USER_COLLECTION is not set in .env file")
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

// func TestPut(t *testing.T) {
// 	collection := setupTestDB(t)
// 	t.Cleanup(func() { cleanupTestDB(t, collection) })
// 	store := NewUserStore(collection)
// 	ctx := context.Background()

// 	testCases := []struct {
// 		name        string
// 		email       string
// 		expectedErr error
// 		isEmpty     bool
// 	}{
// 		{
// 			name:        "should create new user with unique email",
// 			email:       email,
// 			expectedErr: nil,
// 			isEmpty:     false,
// 		},
// 		{
// 			name:        "should reject duplicate email",
// 			email:       email,
// 			expectedErr: fmt.Errorf("user already exists with this email: %s", email),
// 			isEmpty:     true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			userID, err := store.Put(ctx, "John Doe", tc.email, "password123")
// 			assert.Equal(t, err, tc.expectedErr)
// 			assert.Equal(t, tc.isEmpty, len(userID) == 0)
// 		})
// 	}
// }

// func TestGet(t *testing.T) {
// 	collection := setupTestDB(t)
// 	t.Cleanup(func() {
// 		cleanupTestDB(t, collection)
// 	})

// 	store := NewUserStore(collection)

// 	ctx := context.Background()

// 	insertID, err := store.Put(ctx, name, email, password)
// 	assert.NoError(t, err)

// 	t.Run("should return user when supplied valid userId", func(t *testing.T) {
// 		retrievedUser, err := store.Get(ctx, insertID)
// 		assert.NoError(t, err)
// 		assert.Equal(t, email, retrievedUser.Email)
// 	})

// 	t.Run("should return error when supplied non-existing userId", func(t *testing.T) {
// 		nonExistentUserId := "60d5f87b8f8b4c5a5f8b4567"
// 		retrievedUser, err := store.Get(ctx, nonExistentUserId)
// 		assert.EqualError(t, err, "user not found")
// 		assert.Equal(t, User{}, retrievedUser)
// 	})

// 	t.Run("should return error when supplied an id with invalid format", func(t *testing.T) {
// 		retrievedUser, err := store.Get(ctx, "unformattedID")
// 		assert.ErrorContains(t, err, "invalid ID format")
// 		assert.Equal(t, User{}, retrievedUser)
// 	})
// }

// func TestGetByAuth(t *testing.T) {
// 	collection := setupTestDB(t)
// 	t.Cleanup(func() {
// 		cleanupTestDB(t, collection)
// 	})

// 	store := NewUserStore(collection)

// 	ctx := context.Background()

// 	_, err := store.Put(ctx, name, email, password)
// 	assert.NoError(t, err)

// 	t.Run("should return user when supplied correct credentials", func(t *testing.T) {
// 		retrievedUser, err := store.GetByAuth(ctx, email, password)
// 		assert.NoError(t, err)
// 		assert.Equal(t, email, retrievedUser.Email)
// 	})

// 	t.Run("should return error when supplied incorrect credentials", func(t *testing.T) {
// 		retrievedUser, err := store.GetByAuth(ctx, email, "wrongPassword")
// 		assert.ErrorContains(t, err, "incorrect authentication credentials")
// 		assert.Equal(t, User{}, retrievedUser)

// 		retrievedUser, err = store.GetByAuth(ctx, "wrongEmail", password)
// 		assert.EqualError(t, err, "user not found")
// 		assert.Equal(t, User{}, retrievedUser)
// 	})
// }

// func TestGetStyleField(t *testing.T) {
// 	collection := setupTestDB(t)
// 	t.Cleanup(func() {
// 		cleanupTestDB(t, collection)
// 	})

// 	store := NewUserStore(collection)
// 	ctx := context.Background()
// 	userID, err := store.Put(ctx, name, email, password)
// 	assert.NoError(t, err)

// 	t.Run("should return style field when valid style field is provided", func(t *testing.T) {
// 		styleField := SubjectiveStyleField
// 		newStyle := "bold"
// 		err := store.UpdateStyle(ctx, userID, styleField, newStyle)
// 		assert.NoError(t, err)

// 		result, err := store.GetStyleField(ctx, userID, styleField)
// 		assert.NoError(t, err)
// 		assert.Equal(t, newStyle, result)
// 	})

// 	t.Run("should return error when invalid style field is provided", func(t *testing.T) {
// 		_, err := store.GetStyleField(ctx, userID, "invalidStyleField")
// 		assert.EqualError(t, err, "invalid style: invalidStyleField")
// 	})

// 	t.Run("should return error when user is not found", func(t *testing.T) {
// 		_, err := store.GetStyleField(ctx, primitive.NewObjectID().Hex(), SubjectiveStyleField)
// 		assert.EqualError(t, err, "user not found")
// 	})
// }


func TestUpdateStyle(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, collection)
	})

	store := NewUserStore(collection)

	ctx := context.Background()

	providerID, err := store.Put(ctx, name, email, password)
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		styleField    string
		expectedStyle string
		fieldToCheck  func(user *User) string
	}{
		{
			name:          "should return provider with new subjective style",
			styleField:    SubjectiveStyleField,
			expectedStyle: "sampleSubjectiveStyle",
			fieldToCheck:  func(user *User) string { return user.SubjectiveStyle },
		},
		{
			name:          "should return provider with new objective style",
			styleField:    ObjectiveStyleField,
			expectedStyle: "sampleObjectiveStyle",
			fieldToCheck:  func(user *User) string { return user.ObjectiveStyle },
		},
		{
			name:          "should return provider with new assessment style",
			styleField:    AssessmentAndPlanStyleField,
			expectedStyle: "sampleAssessmentStyle",
			fieldToCheck:  func(user *User) string { return user.AssessmentAndPlanStyle },
		},
		{
			name:          "should return provider with new summary style",
			styleField:    SummaryStyleField,
			expectedStyle: "sampleSummaryStyle",
			fieldToCheck:  func(user *User) string { return user.SummaryStyle },
		},
		{
			name:          "should return provider with new patientInstructions style",
			styleField:    PatientInstructionsStyleField,
			expectedStyle: "samplePatientInstructionsStyle",
			fieldToCheck:  func(user *User) string { return user.PatientInstructionsStyle },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.UpdateStyle(ctx, providerID, tc.styleField, tc.expectedStyle)
			assert.NoError(t, err)

			provider, err := store.Get(ctx, providerID)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStyle, tc.fieldToCheck(&provider))
		})
	}

	t.Run("should return error when invalid style is provided", func(t *testing.T) {
		expectedStyle := "sampleStyle"

		err := store.UpdateStyle(ctx, providerID, SubjectiveStyleField, expectedStyle)
		assert.NoError(t, err, "correct field was updated")

		err = store.UpdateStyle(ctx, providerID, "InvalidStyle", "incorrectStyleUpdate")
		assert.Error(t, err, "invalid field should return error")

		provider, err := store.Get(ctx, providerID)
		assert.NoError(t, err)
		assert.Equal(t, expectedStyle, provider.SubjectiveStyle)
	})

	t.Run("should return error when document doesn't exist", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		err := store.UpdateStyle(ctx, nonExistentID, SubjectiveStyleField, "someStyle")
		assert.Error(t, err, "non-existent ID should return error")
		assert.Contains(t, err.Error(), "no document found", "Error message should contain 'no document found'")
	})

	t.Run("should return error when ObjectID is invalid", func(t *testing.T) {
		invalidObjectID := "invalid-object-id"
		err := store.UpdateStyle(ctx, invalidObjectID, SubjectiveStyleField, "someStyle")
		assert.Error(t, err, "invalid ObjectID should return error")
		assert.Contains(t, err.Error(), "invalid ID format", "Error message should contain 'invalid ID format'")

	})
}
