package reports

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

type MockReportsStore struct {
	mock.Mock
}

func (m *MockReportsStore) Put(ctx context.Context, name, providerID string, timestamp time.Time, duration float64, isFollowUp bool, pronouns, lastVisitID string, usedDiarization bool) (string, error) {
	args := m.Called(ctx, name, providerID, timestamp, duration, isFollowUp, pronouns)
	return args.String(0), args.Error(1)
}

func (m *MockReportsStore) Get(ctx context.Context, reportId string) (Report, error) {
	args := m.Called(ctx, reportId)
	return args.Get(0).(Report), args.Error(1)
}

func (m *MockReportsStore) UpdateStatus(ctx context.Context, reportId string, status string) error {
	args := m.Called(ctx, reportId, status)
	return args.Error(0)
}

func (m *MockReportsStore) GetTranscription(ctx context.Context, reportId string) (RetrievedReportTranscripts, error) {
	args := m.Called(ctx, reportId)
	return args.Get(0).(RetrievedReportTranscripts), args.Error(1)
}

func (m *MockReportsStore) GetAll(ctx context.Context, userId string) ([]Report, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]Report), args.Error(1)
}

func (m *MockReportsStore) Delete(ctx context.Context, reportId string) error {
	args := m.Called(ctx, reportId)
	return args.Error(0)
}

func (m *MockReportsStore) MarkRead(ctx context.Context, reportId string) error {
	args := m.Called(ctx, reportId)
	return args.Error(0)
}

func (m *MockReportsStore) MarkUnread(ctx context.Context, reportId string) error {
	args := m.Called(ctx, reportId)
	return args.Error(0)
}

func (m *MockReportsStore) UpdateReport(ctx context.Context, reportId string, batchedUpdates bson.D) error {
	args := m.Called(ctx, reportId, batchedUpdates)
	return args.Error(0)
}

func (m *MockReportsStore) Validate(report *Report) error {
	args := m.Called(report)
	return args.Error(0)
}
