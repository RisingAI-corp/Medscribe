package inferenceService

import (
	inferencestore "Medscribe/inference/store"
	"Medscribe/reports"
	transcriber "Medscribe/transcription"
	"Medscribe/user"
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Constants for Test Data ---
const (
	reportID          = "report-123"
	providerID        = "provider-123"
	sampleContentData = "test query response data"
)

func sortedBsonD(doc bson.D) bson.D {
	sorted := make(bson.D, len(doc))
	copy(sorted, doc)

	for i, elem := range sorted {
		if inner, ok := elem.Value.(bson.D); ok {
			sorted[i].Value = sortedBsonD(inner)
		}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	return sorted
}

type testEnv struct {
	service        InferenceService
	reportCfg      *ReportRequest
	contentChan    chan ContentChanPayload
	transcriber    *transcriber.MockTranscription
	reports        *reports.MockReports
	inferenceStore *inferencestore.MockInferenceStore
	users          *user.MockUserStore
}

func setupTestEnvironment(t *testing.T) *testEnv {
	t.Helper()
	transcriber := new(transcriber.MockTranscription)
	reportStore := new(reports.MockReports)
	inferenceStore := new(inferencestore.MockInferenceStore)
	userStore := new(user.MockUserStore)

	service := NewInferenceService(reportStore, transcriber, inferenceStore, userStore)

	reportCfg := &ReportRequest{
		AudioBytes:       []byte("dummy audio"),
		TranscribedAudio: "",
		ReportContents: []ReportContentSection{
			{ContentType: reports.Subjective, Content: reports.Content},
			{ContentType: reports.Objective, Content: reports.Content},
			{ContentType: reports.Assessment, Content: reports.Content},
			{ContentType: reports.Planning, Content: reports.Content},
			{ContentType: reports.Summary, Content: reports.Content},
		},
		PatientName: "John Doe",
		ProviderID:  "provider-123",
		Timestamp:   time.Now(),
		Duration:    60,
	}

	contentChan := make(chan ContentChanPayload, 10)

	return &testEnv{
		service:        service,
		reportCfg:      reportCfg,
		contentChan:    contentChan,
		transcriber:    transcriber,
		reports:        reportStore,
		inferenceStore: inferenceStore,
		users:          userStore,
	}
}

func TestSuccessfulReportGeneration(t *testing.T) {
	env := setupTestEnvironment(t)

	env.transcriber.
		On("Transcribe", mock.Anything, []byte("dummy audio")).
		Return("dummy transcript", nil)

	env.reports.
		On("Put", mock.Anything, "John Doe", "provider-123", mock.Anything, float64(60), false, mock.Anything).
		Return(reportID, nil)

	env.inferenceStore.
		On("Query", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Return("test query response data", nil)

	env.reports.
		On("UpdateReport", mock.Anything, reportID, mock.MatchedBy(func(batchedUpdates bson.D) bool {
			expectedBatchedUpdates := bson.D{
				{Key: reports.Subjective, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.Objective, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.Assessment, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.Planning, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.Summary, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.FinishedGenerating, Value: true},
			}
			sortedBatchUpdates := sortedBsonD(batchedUpdates)
			sortedExpectedBatchedUpdates := sortedBsonD(expectedBatchedUpdates)
			return reflect.DeepEqual(sortedExpectedBatchedUpdates, sortedBatchUpdates)
		})).
		Return(nil)

	err := env.service.GenerateReportPipeline(context.Background(), env.reportCfg, env.contentChan)
	require.NoError(t, err)

	received := make(map[string]interface{})

	for msg := range env.contentChan {
		received[msg.Key] = msg.Value
	}

	expected := map[string]interface{}{
		"_id":                      reportID,
		reports.Subjective:         sampleContentData,
		reports.Objective:          sampleContentData,
		reports.Assessment:         sampleContentData,
		reports.Planning:           sampleContentData,
		reports.Summary:            sampleContentData,
		reports.FinishedGenerating: true,
	}

	require.Equal(t, expected, received)

	env.transcriber.AssertExpectations(t)
	env.reports.AssertExpectations(t)
	env.inferenceStore.AssertExpectations(t)

}

func TestFailedReportGeneration_InvalidContentType(t *testing.T) {
	env := setupTestEnvironment(t)
	env.reportCfg.ReportContents = []ReportContentSection{
		{ContentType: "invalid_type", Content: reports.Content},
	}

	err := env.service.GenerateReportPipeline(context.Background(), env.reportCfg, env.contentChan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not a valid contentType")

	env.reports.AssertNotCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	env.reports.AssertNotCalled(t, "UpdateReport", mock.Anything, mock.Anything, mock.Anything)
	env.inferenceStore.AssertNotCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything)
}

func TestFailedReportGeneration_Transcription(t *testing.T) {
	env := setupTestEnvironment(t)

	env.transcriber.
		On("Transcribe", mock.Anything, []byte("dummy audio")).
		Return("", errors.New("transcription error"))

	err := env.service.GenerateReportPipeline(context.Background(), env.reportCfg, env.contentChan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "transcribing audio")

	env.transcriber.AssertExpectations(t)
	env.reports.AssertNotCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	env.reports.AssertNotCalled(t, "UpdateReport", mock.Anything, mock.Anything, mock.Anything)
	env.inferenceStore.AssertNotCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything)
}

// TestFailedReportGeneration_Put triggers a failure in the Put method.
func TestFailedReportGeneration_Put(t *testing.T) {
	env := setupTestEnvironment(t)

	env.transcriber.
		On("Transcribe", mock.Anything, []byte("dummy audio")).
		Return("dummy transcript", nil)

	env.reports.
		On("Put", mock.Anything, "John Doe", "provider-123", mock.Anything, float64(60), false, mock.Anything).
		Return("", errors.New("put error"))

	err := env.service.GenerateReportPipeline(context.Background(), env.reportCfg, env.contentChan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error storing report")

	env.transcriber.AssertExpectations(t)
	env.reports.AssertExpectations(t)
	env.reports.AssertNotCalled(t, "UpdateReport", mock.Anything, mock.Anything, mock.Anything)
	env.inferenceStore.AssertNotCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything)
}

// TestFailedReportGeneration_UpdateReport triggers a failure in the UpdateReport method.
func TestFailedReportGeneration_UpdateReport(t *testing.T) {
	env := setupTestEnvironment(t)

	env.transcriber.
		On("Transcribe", mock.Anything, []byte("dummy audio")).
		Return("dummy transcript", nil)

	env.reports.
		On("Put", mock.Anything, "John Doe", "provider-123", mock.Anything, float64(60), false, mock.Anything).
		Return(reportID, nil)

	env.inferenceStore.
		On("Query", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Return("generated section content", nil)

	env.reports.
		On("UpdateReport", mock.Anything, reportID, mock.Anything).
		Return(errors.New("update error"))

	err := env.service.GenerateReportPipeline(context.Background(), env.reportCfg, env.contentChan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "updating report")

	env.transcriber.AssertExpectations(t)
	env.reports.AssertExpectations(t)
	env.inferenceStore.AssertExpectations(t)
}

func TestRegenerateReport_Valid(t *testing.T) {
	env := setupTestEnvironment(t)

	validReportContents := []ReportContentSection{
		{ContentType: reports.Subjective, Content: reports.ContentData},
		{ContentType: reports.Objective, Content: reports.Content},
		{ContentType: reports.Assessment, Content: reports.Content},
		{ContentType: reports.Planning, Content: reports.Content},
		{ContentType: reports.Summary, Content: reports.Content},
	}

	updates := bson.D{
		{Key: reports.Pronouns, Value: reports.HE},
	}

	env.inferenceStore.
		On("Query", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Return("test query response data", nil)

	env.reports.
		On("UpdateReport", mock.Anything, reportID, bson.D{
			{Key: "pronouns", Value: reports.HE},
			{Key: reports.FinishedGenerating, Value: false},
		}).
		Return(nil)

	env.reports.
		On("UpdateReport", mock.Anything, reportID, mock.MatchedBy(func(batchedUpdates bson.D) bool {
			expectedBatchedUpdates := bson.D{
				{Key: reports.Subjective, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.Objective, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.Assessment, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.Planning, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.Summary, Value: bson.D{{Key: reports.ContentData, Value: sampleContentData}, {Key: reports.Loading, Value: false}}},
				{Key: reports.FinishedGenerating, Value: true},
			}

			sortedBatchUpdates := sortedBsonD(batchedUpdates)
			sortedExpectedBatchedUpdates := sortedBsonD(expectedBatchedUpdates)
			return reflect.DeepEqual(sortedExpectedBatchedUpdates, sortedBatchUpdates)
		})).
		Return(nil)

	reportRequest := &ReportRequest{
		ID:                reportID,
		ReportContents:    validReportContents,
		Updates:           updates,
		SubjectiveContent: "here is sample subjective content",
	}
	err := env.service.RegenerateReport(context.Background(), env.contentChan, reportRequest)
	require.NoError(t, err)

	received := make(map[string]interface{})
	for msg := range env.contentChan {
		received[msg.Key] = msg.Value
	}

	expected := map[string]interface{}{
		reports.Subjective:         sampleContentData,
		reports.Objective:          sampleContentData,
		reports.Assessment:         sampleContentData,
		reports.Planning:           sampleContentData,
		reports.Summary:            sampleContentData,
		reports.FinishedGenerating: true,
	}
	require.Equal(t, expected, received)

	env.reports.AssertExpectations(t)
	env.inferenceStore.AssertExpectations(t)
}

func TestRegenerateReport_InvalidInputs(t *testing.T) {
	t.Run("should return error when report content type is invalid", func(t *testing.T) {
		env := setupTestEnvironment(t)
		invalidReportContents := []ReportContentSection{
			{
				ContentType: "InvalidContentType",
				Content:     "invalid content",
			},
		}
		reportRequest := &ReportRequest{
			ID:             reportID,
			ReportContents: invalidReportContents,
			Updates:        bson.D{},
		}
		err := env.service.RegenerateReport(context.Background(), env.contentChan, reportRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "is not a valid contentType")

		env.reports.AssertNotCalled(t, "UpdateReport", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return error when no updates are provided", func(t *testing.T) {
		env := setupTestEnvironment(t)
		reportRequest := &ReportRequest{
			ID:             reportID,
			ReportContents: []ReportContentSection{{ContentType: reports.Subjective, Content: reports.Content}},
		}
		err := env.service.RegenerateReport(context.Background(), env.contentChan, reportRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no updates provided")

		env.reports.AssertNotCalled(t, "UpdateReport", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return error when inference store query fails", func(t *testing.T) {
		env := setupTestEnvironment(t)

		env.inferenceStore.
			On("Query", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).
			Return("", errors.New("inference error"))

		env.reports.
			On("UpdateReport", mock.Anything, reportID, bson.D{{Key: reports.FinishedGenerating, Value: false}}).
			Return(nil).Once()

		validReportContents := []ReportContentSection{
			{ContentType: reports.Subjective, Content: reports.Content},
			{ContentType: reports.Objective, Content: reports.Content},
			{ContentType: reports.Assessment, Content: reports.Content},
			{ContentType: reports.Planning, Content: reports.Content},
			{ContentType: reports.Summary, Content: reports.Content},
		}
		reportRequest := &ReportRequest{
			ID:             reportID,
			ReportContents: validReportContents,
			Updates:        bson.D{},
		}

		err := env.service.RegenerateReport(context.Background(), env.contentChan, reportRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "error generating report sections")

		env.reports.AssertNumberOfCalls(t, "UpdateReport", 1)
	})

	t.Run("should return error when setting loading to true fails", func(t *testing.T) {
		env := setupTestEnvironment(t)
		env.reports.
			On("UpdateReport", context.Background(), reportID, bson.D{{Key: reports.FinishedGenerating, Value: false}}).
			Return(errors.New("update error"))

		reportRequest := &ReportRequest{
			ID:             reportID,
			ReportContents: []ReportContentSection{{ContentType: reports.Subjective, Content: reports.Content}},
			Updates:        bson.D{},
		}

		err := env.service.RegenerateReport(context.Background(), env.contentChan, reportRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "error updating loading status before report regeneration:")

		env.reports.AssertExpectations(t)
	})

	t.Run("should return error when UpdateReport fails after sending batch updates", func(t *testing.T) {
		env := setupTestEnvironment(t)

		env.inferenceStore.
			On("Query", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).
			Return(sampleContentData, nil)

		env.reports.
			On("UpdateReport", context.Background(), reportID, bson.D{{Key: reports.FinishedGenerating, Value: false}}).
			Return(nil).Once()

		expectedUpdate := primitive.D{
			primitive.E{Key: "summary", Value: primitive.D{primitive.E{Key: "data", Value: "test query response data"}, primitive.E{Key: "loading", Value: false}}},
			primitive.E{Key: "subjective", Value: primitive.D{primitive.E{Key: "data", Value: "test query response data"}, primitive.E{Key: "loading", Value: false}}},
			primitive.E{Key: "assessment", Value: primitive.D{primitive.E{Key: "data", Value: "test query response data"}, primitive.E{Key: "loading", Value: false}}},
			primitive.E{Key: "objective", Value: primitive.D{primitive.E{Key: "data", Value: "test query response data"}, primitive.E{Key: "loading", Value: false}}},
			primitive.E{Key: "planning", Value: primitive.D{primitive.E{Key: "data", Value: "test query response data"}, primitive.E{Key: "loading", Value: false}}},
			primitive.E{Key: "finishedGenerating", Value: true},
		}
		env.reports.
			On("UpdateReport", mock.Anything, "report-123", mock.MatchedBy(func(arg primitive.D) bool {
				return len(arg) == len(expectedUpdate)
			})).Return(errors.New("update error")).Once()

		reportRequest := &ReportRequest{
			ID:             reportID,
			ReportContents: []ReportContentSection{{ContentType: reports.Subjective, Content: reports.Content}},
			Updates:        bson.D{},
		}

		err := env.service.RegenerateReport(context.Background(), env.contentChan, reportRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "error updating report after regeneration:")

		env.reports.AssertExpectations(t)
	})

}

func TestLearnStyle_Valid(t *testing.T) {
	env := setupTestEnvironment(t)

	contentSection := reports.Subjective
	content := "some valid content"
	newStyle := "learned style"

	env.inferenceStore.
		On("Query", mock.Anything, mock.Anything, 100).
		Return(newStyle, nil)

	// userStore.UpdateStyle()
	env.users.On("UpdateStyle", context.Background(), providerID, user.SubjectiveStyleField, newStyle).Return(nil)

	err := env.service.LearnStyle(context.Background(), providerID, contentSection, content)
	require.NoError(t, err)

	env.inferenceStore.AssertExpectations(t)
	env.reports.AssertExpectations(t)
	env.users.AssertExpectations(t)
}

func TestLearnStyle_Invalid(t *testing.T) {
	t.Run("should return error when report content type is invalid", func(t *testing.T) {
		env := setupTestEnvironment(t)

		contentSection := "invalid content type"
		content := "some valid content"

		err := env.service.LearnStyle(context.Background(), providerID, contentSection, content)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid content section")

		env.inferenceStore.AssertNotCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything)
		env.users.AssertNotCalled(t, "UpdateStyle", mock.Anything, mock.Anything, mock.Anything)

	})

	t.Run("should return error when content is empty", func(t *testing.T) {
		env := setupTestEnvironment(t)

		err := env.service.LearnStyle(context.Background(), providerID, "", "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot learn from empty content")

		env.inferenceStore.AssertNotCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything)
		env.users.AssertNotCalled(t, "UpdateStyle", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return error when inference store query fails", func(t *testing.T) {
		env := setupTestEnvironment(t)

		contentSection := reports.Subjective
		content := "some valid content"

		env.inferenceStore.
			On("Query", mock.Anything, mock.Anything, 100).
			Return("", errors.New("query error"))

		err := env.service.LearnStyle(context.Background(), providerID, contentSection, content)
		require.Error(t, err)
		require.Contains(t, err.Error(), "error querying for style")

		env.reports.AssertNotCalled(t, "UpdateReport", mock.Anything, mock.Anything, mock.Anything)
		env.inferenceStore.AssertExpectations(t)
		env.users.AssertNotCalled(t, "UpdateStyle", mock.Anything, mock.Anything, mock.Anything)

	})

	t.Run("should return error when Update user styles fails", func(t *testing.T) {
		env := setupTestEnvironment(t)

		contentSection := reports.Subjective
		content := "some valid content"

		newStyle := "learned style"

		env.inferenceStore.
			On("Query", mock.Anything, mock.Anything, 100).
			Return(newStyle, nil)

		env.users.On("UpdateStyle", context.Background(), providerID, user.SubjectiveStyleField, newStyle).Return(errors.New("update error"))

		err := env.service.LearnStyle(context.Background(), providerID, contentSection, content)
		require.Error(t, err)
		require.Contains(t, err.Error(), "error updating style")

		env.inferenceStore.AssertExpectations(t)
		env.users.AssertExpectations(t)
	})

}
