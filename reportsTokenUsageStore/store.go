package reportsTokenUsage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TokenUsageEntry tracks the number of tokens used per section of a report
type TokenUsageEntry struct {
	ReportID    primitive.ObjectID `bson:"reportId"`
	ProviderID  string             `bson:"providerId"`
	Timestamp   primitive.DateTime `bson:"timestamp"`
	TokenUsage  map[string]int     `bson:"tokenUsage"`
	TotalTokens int                `bson:"totalTokens"`
}

type TokenUsageStore interface {
	Insert(ctx context.Context, entry TokenUsageEntry) error
	UpdateSectionTokens(ctx context.Context, reportId primitive.ObjectID, section string, tokens int) error
	GetByReportID(ctx context.Context, reportId primitive.ObjectID) (TokenUsageEntry, error)
}

type tokenUsageStore struct {
	collection *mongo.Collection
}

func NewTokenUsageStore(collection *mongo.Collection) TokenUsageStore {
	return &tokenUsageStore{collection: collection}
}
func (s *tokenUsageStore) Insert(ctx context.Context, entry TokenUsageEntry) error {
	if entry.ReportID.IsZero() {
		return errors.New("reportId cannot be empty")
	}
	if entry.ProviderID == "" {
		return errors.New("providerId cannot be empty")
	}
	if entry.Timestamp.Time().IsZero() {
		entry.Timestamp = primitive.NewDateTimeFromTime(time.Now())
	}
	if entry.TotalTokens == 0 {
		tokens := 0
		for _, v := range entry.TokenUsage {
			tokens += v
		}
		entry.TotalTokens = tokens
	}
	_, err := s.collection.InsertOne(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to insert token usage entry: %v", err)
	}
	return nil
}

func (s *tokenUsageStore) UpdateSectionTokens(ctx context.Context, reportId primitive.ObjectID, section string, tokens int) error {
	if reportId.IsZero() || section == "" || tokens < 0 {
		return errors.New("invalid input for updating token usage")
	}

	filter := bson.M{"reportId": reportId}
	update := bson.M{
		"$set":         bson.M{fmt.Sprintf("tokenUsage.%s", section): tokens},
		"$setOnInsert": bson.M{"timestamp": primitive.NewDateTimeFromTime(time.Now())},
	}
	opts := options.Update().SetUpsert(true)

	_, err := s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update token usage: %v", err)
	}
	return nil
}

func (s *tokenUsageStore) GetByReportID(ctx context.Context, reportId primitive.ObjectID) (TokenUsageEntry, error) {
	if reportId.IsZero() {
		return TokenUsageEntry{}, errors.New("invalid reportId")
	}

	var result TokenUsageEntry
	err := s.collection.FindOne(ctx, bson.M{"reportId": reportId}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return TokenUsageEntry{}, fmt.Errorf("no token usage entry found for reportId: %s", reportId.Hex())
		}
		return TokenUsageEntry{}, fmt.Errorf("error retrieving token usage: %v", err)
	}

	return result, nil
}
