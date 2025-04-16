package reportsTokenUsage

import (
	"context"
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
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("failed to load env: %v", err)
	}
}

func setupTestCollection(t *testing.T) *mongo.Collection {
	t.Helper()
	loadEnv(t)

	uri := os.Getenv("MONGODB_URI_DEV")
	db := os.Getenv("MONGODB_DB")
	collection := os.Getenv("MONGODB_REPORTS_TOKEN_USAGE_STORE")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("mongo connection failed: %v", err)
	}

	return client.Database(db).Collection(collection)
}

func cleanupCollection(t *testing.T, collection *mongo.Collection) {
	t.Helper()
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}

func TestTokenUsageStore_InsertAndGet(t *testing.T) {
	collection := setupTestCollection(t)
	t.Cleanup(func() { cleanupCollection(t, collection) })

	store := NewTokenUsageStore(collection)
	ctx := context.Background()

	reportID := primitive.NewObjectID()
	entry := TokenUsageEntry{
		ReportID:   reportID,
		ProviderID: "provider-001",
		Timestamp:  primitive.NewDateTimeFromTime(time.Now()),
		TokenUsage: map[string]int{
			"subjective":            100,
			"objective":             120,
			"assessmentAndPlanning": 130,
			"condensedSummary":      80,
			"sessionSummary":        90,
			"patientInstructions":   110,
		},
	}

	err := store.Insert(ctx, entry)
	assert.NoError(t, err)

	retrieved, err := store.GetByReportID(ctx, reportID)
	assert.NoError(t, err)
	assert.Equal(t, entry.ReportID, retrieved.ReportID)
	assert.Equal(t, entry.ProviderID, retrieved.ProviderID)
	assert.Equal(t, entry.TokenUsage["objective"], retrieved.TokenUsage["objective"])
}

func TestTokenUsageStore_UpdateSectionTokens(t *testing.T) {
	collection := setupTestCollection(t)
	t.Cleanup(func() { cleanupCollection(t, collection) })

	store := NewTokenUsageStore(collection)
	ctx := context.Background()

	reportID := primitive.NewObjectID()
	initial := TokenUsageEntry{
		ReportID:   reportID,
		ProviderID: "provider-002",
		Timestamp:  primitive.NewDateTimeFromTime(time.Now()),
		TokenUsage: map[string]int{},
	}

	err := store.Insert(ctx, initial)
	assert.NoError(t, err)

	err = store.UpdateSectionTokens(ctx, reportID, "sessionSummary", 456)
	assert.NoError(t, err)

	updated, err := store.GetByReportID(ctx, reportID)
	assert.NoError(t, err)
	assert.Equal(t, 456, updated.TokenUsage["sessionSummary"])
}
