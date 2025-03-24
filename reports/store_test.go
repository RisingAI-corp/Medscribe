package reports

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in .env file")
	}

	dbName := os.Getenv("MONGODB_DB")
	if mongoURI == "" {
		log.Fatal("MONGODB_DB is not set in .env file")
	}

	collectionName := os.Getenv("MONGODB_REPORT_COLLECTION_DEV")
	if mongoURI == "" {
		log.Fatal("MONGODB_REPORT_COLLECTION is not set in .env file")
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

func TestPut(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })

	store := NewReportsStore(collection)
	ctx := context.Background()

	testCases := []struct {
		name          string
		providerID    string
		duration      float64
		pronouns      string
		expectedErr   error
		isResultEmpty bool
		description   string
	}{
		{
			name:          "sampleReport",
			providerID:    "providerID123",
			expectedErr:   nil,
			isResultEmpty: false,
			description:   "should return reportId when supplied with correct inputs",
			duration:      1,
			pronouns:      HE,
		},
		{
			name:          "",
			providerID:    "",
			expectedErr:   errors.New("name cannot be an empty string"),
			isResultEmpty: true,
			description:   "should throw error when name is an empty string",
		},
		{
			name:          "sampleReport",
			providerID:    "",
			expectedErr:   errors.New("providerId cannot be an empty string"),
			isResultEmpty: true,
			description:   "should throw error when providerId is an empty string",
		},
		{
			name:          "invalid duration",
			providerID:    "providerID123",
			duration:      0, // Invalid duration
			pronouns:      HE,
			expectedErr:   errors.New("duration must be greater than 0"),
			isResultEmpty: true,
			description:   "should throw error when duration is less than or equal to 0",
		},
		{
			name:          "invalid pronouns",
			providerID:    "providerID123",
			duration:      60,
			pronouns:      "INVALID",
			expectedErr:   fmt.Errorf("pronouns must be either '%s', '%s', or '%s'", HE, SHE, THEY),
			isResultEmpty: true,
			description:   "should throw error when pronouns is invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			reportId, err := store.Put(ctx, tc.name, tc.providerID, time.Now(), tc.duration, false, tc.pronouns)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.isResultEmpty, reportId == "")
		})
	}
}

func TestGet(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })

	store := NewReportsStore(collection)

	ctx := context.Background()

	name := "John Doe"
	providerId := "provider123"

	reportID, err := store.Put(ctx, name, providerId, time.Now(), 1, false, HE)
	assert.NoError(t, err)

	t.Run("should fetch report when ID exists", func(t *testing.T) {
		report, err := store.Get(ctx, reportID)
		assert.NoError(t, err)
		assert.Equal(t, report.ID.Hex(), reportID)
	})

	t.Run("should error when ID is invalid", func(t *testing.T) {
		report, err := store.Get(ctx, "invalidReportID")
		assert.ErrorContains(t, err, "invalid ID format:")
		assert.Zero(t, report)
	})
}

func TestGetAll(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })

	store := NewReportsStore(collection)
	ctx := context.Background()

	providerID1 := "provider123"
	providerID2 := "provider456"
	providerID3 := "provider789"

	reportID1, err := store.Put(ctx, "Report1", providerID1, time.Now(), 1, false, HE)
	assert.NoError(t, err)
	reportID2, err := store.Put(ctx, "Report2", providerID1, time.Now(), 1, false, HE)
	assert.NoError(t, err)

	reportID3, err := store.Put(ctx, "Report3", providerID2, time.Now(), 1, false, HE)
	assert.NoError(t, err)

	reportID4, err := store.Put(ctx, "Report4", providerID3, time.Now(), 1, false, HE)
	assert.NoError(t, err)

	testCases := []struct {
		providerID    string
		expectedIDs   []string
		expectedError error
		description   string
	}{
		{
			providerID:    providerID1,
			expectedIDs:   []string{reportID1, reportID2},
			expectedError: nil,
			description:   "should fetch all reports for providerID1",
		},
		{
			providerID:    providerID2,
			expectedIDs:   []string{reportID3},
			expectedError: nil,
			description:   "should fetch all reports for providerID2",
		},
		{
			providerID:    providerID3,
			expectedIDs:   []string{reportID4},
			expectedError: nil,
			description:   "should fetch all reports for providerID3",
		},
		{
			providerID:    "",
			expectedIDs:   nil,
			expectedError: errors.New("missing provider ID"),
			description:   "should return error when providerID is empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			reports, err := store.GetAll(ctx, tc.providerID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, len(tc.expectedIDs), len(reports))
		})
	}
}

func TestDelete(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })

	store := NewReportsStore(collection)
	ctx := context.Background()

	t.Run("should successfully delete existing report", func(t *testing.T) {
		reportID, err := store.Put(ctx, "Test Report", "provider123", time.Now(), 120, false, HE)
		assert.NoError(t, err, "failed to create test report")

		err = store.Delete(ctx, reportID)
		assert.NoError(t, err, "failed to delete report")

		_, err = store.Get(ctx, reportID)
		assert.Error(t, err, "report should not exist after deletion")
	})

	t.Run("should return error when deleting non-existent report", func(t *testing.T) {
		err := store.Delete(ctx, primitive.NewObjectID().Hex())
		assert.Error(t, err, "deleting non-existent report should return error")
		assert.Contains(t, err.Error(), "report not found")
	})

	t.Run("should return error with invalid report ID format", func(t *testing.T) {
		err := store.Delete(ctx, "invalid-id")
		assert.Error(t, err, "invalid ID format should return error")
		assert.Contains(t, err.Error(), "invalid ID format")
	})
}

func TestGetTranscription(t *testing.T) {
	collection := setupTestDB(t)
	// t.Cleanup(func() { cleanupTestDB(t, collection) })

	store := NewReportsStore(collection)
	ctx := context.Background()

	t.Run("should return transcription for existing report", func(t *testing.T) {
		reportID, err := store.Put(ctx, "Test Report", "provider123", time.Now(), 120, false, HE)
		assert.NoError(t, err, "failed to create test report")

		providerID, transcription, err := store.GetTranscription(ctx, reportID)
		assert.NoError(t, err)
		assert.Empty(t, transcription)
		assert.Equal(t, providerID, "provider123")
	})

	t.Run("should return error when getting transcription for non-existent report", func(t *testing.T) {
		providerID, transcription, err := store.GetTranscription(ctx, primitive.NewObjectID().Hex())
		assert.Error(t, err, "getting transcription for non-existent report should return error")
		assert.Contains(t, err.Error(), "failed to retrieve transcript")
		assert.Empty(t, transcription)
		assert.Empty(t, providerID)
	})
}


func TestUpdateReport_SuccessfulUpdates(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })

	store := NewReportsStore(collection)
	ctx := context.Background()

	reportID, err := store.Put(ctx, "Test Report", "provider123", time.Now(), 120, false, HE)
	assert.NoError(t, err)

	validUpdates := bson.D{
		{Key: "name", Value: "Updated Name"},
		{Key: "duration", Value: 180},
		{Key: "pronouns", Value: SHE},
		{Key: "subjective", Value: bson.D{{Key: "data", Value: "test data"}}},
	}

	err = store.UpdateReport(ctx, reportID, validUpdates)
	assert.NoError(t, err)

	updatedReport, err := store.Get(ctx, reportID)
	assert.NoError(t, err)

	assert.Equal(t, "Updated Name", updatedReport.Name)
	assert.Equal(t, float64(180), updatedReport.Duration)
	assert.Equal(t, SHE, updatedReport.Pronouns)
	assert.Equal(t, "test data", updatedReport.Subjective.Data)
}

func TestUpdateReport_ValidationFailures(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })

	store := NewReportsStore(collection)
	ctx := context.Background()

	reportID, err := store.Put(ctx, "Test Report", "provider123", time.Now(), 120, false, HE)
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		updates     bson.D
		expectedErr string
		description string
	}{
		{
			name: "Empty ProviderID",
			updates: bson.D{
				{Key: "providerID", Value: ""},
			},
			expectedErr: "ProviderID cannot be empty",
			description: "should fail when ProviderID is empty",
		},
		{
			name: "Empty Name",
			updates: bson.D{
				{Key: "name", Value: ""},
			},
			expectedErr: "name cannot be empty",
			description: "should fail when Name is empty",
		},
		{
			name: "Invalid Duration",
			updates: bson.D{
				{Key: "duration", Value: -1},
			},
			expectedErr: "duration must be greater than 0",
			description: "should fail when duration is less than 0",
		},
		{
			name: "Invalid Pronouns",
			updates: bson.D{
				{Key: "pronouns", Value: "INVALID"},
			},
			expectedErr: fmt.Sprintf("pronouns must be either '%s', '%s', or '%s'", HE, SHE, THEY),
			description: "should fail when Pronouns is invalid",
		},
		{
			name: "Invalid PatientOrClient",
			updates: bson.D{
				{Key: "patientOrClient", Value: "INVALID"},
			},
			expectedErr: fmt.Sprintf("PatientOrClient must be either '%s' or '%s'", Patient, Client),
			description: "should fail when PatientOrClient is invalid",
		},
		{
			name: "Empty SessionSummary",
			updates: bson.D{
				{Key: "sessionSummary", Value: ""},
			},
			expectedErr: "sessionSummary cannot be empty",
			description: "should fail when sessionSummary is empty",
		},
		{
			name: "Empty CondensedSummary",
			updates: bson.D{
				{Key: "condensedSummary", Value: ""},
			},
			expectedErr: "condensedSummary cannot be empty",
			description: "should fail when condensedSummary is empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := store.UpdateReport(ctx, reportID, tc.updates)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestUpdateReport_InvalidUpdates(t *testing.T) {
	collection := setupTestDB(t)
	t.Cleanup(func() { cleanupTestDB(t, collection) })

	store := NewReportsStore(collection)
	ctx := context.Background()

	reportID, err := store.Put(ctx, "Test Report", "provider123", time.Now(), 120, false, HE)
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		updates     bson.D
		expectedErr string
		description string
	}{
		{
			name: "Invalid Field Type",
			updates: bson.D{
				{Key: "duration", Value: "INVALID"},
			},
			expectedErr: "type mismatch for key 'duration': expected 'float64', got 'string'",
			description: "should fail when updates contain invalid field types",
		},
		{
			name: "Unknown Field",
			updates: bson.D{
				{Key: "unknownfield", Value: "value"},
			},
			expectedErr: fmt.Sprintf("key '%s' not found in reportMap", "unknownfield"),
			description: "should fail when updates contain unknown fields",
		},
		{
			name: "invalid nested field Field",
			updates: bson.D{
				{Key: "subjective", Value: bson.D{{Key: "nonexistentfield", Value: "test data"}}},
			},
			expectedErr: fmt.Sprintf("key '%s' not found in reportMap", "nonexistentfield"),
			description: "should fail when updates for a nested fields that doesn't exist",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := store.UpdateReport(ctx, reportID, tc.updates)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}
