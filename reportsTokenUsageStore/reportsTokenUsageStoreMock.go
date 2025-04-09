package reportsTokenUsage

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockTokenUsageStore struct {
	mock.Mock
}


func (m *MockTokenUsageStore) Insert(ctx context.Context, entry TokenUsageEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockTokenUsageStore) UpdateSectionTokens(ctx context.Context, reportID primitive.ObjectID, section string, tokens int) error {
	args := m.Called(ctx, reportID, section, tokens)
	return args.Error(0)
}

func (m *MockTokenUsageStore) GetByReportID(ctx context.Context, reportID primitive.ObjectID) (TokenUsageEntry, error) {
	args := m.Called(ctx, reportID)
	return args.Get(0).(TokenUsageEntry), args.Error(1)
}