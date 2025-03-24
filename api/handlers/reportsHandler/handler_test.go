package reportsHandler

import (
	"Medscribe/api/middleware"
	inferenceService "Medscribe/inference/service"
	"Medscribe/reports"
	"Medscribe/user"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// Constants for test data
const (
	testUserID    = "test123"
	testReportID  = "report123"
	testAudioData = "test audio data"
)

func TestGenerateReport(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)
	t.Run("should generate report when request is valid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add metadata
		metadata := map[string]interface{}{
			"providerID": testUserID,
			"name":       "Test Report",
			"timestamp":  time.Now(),
			"duration":   30,
		}
		metadataBytes, err := json.Marshal(metadata)
		require.NoError(t, err)

		err = writer.WriteField("metadata", string(metadataBytes))
		require.NoError(t, err)

		// Add audio file
		part, err := writer.CreateFormFile("audio", "test.wav")
		require.NoError(t, err)
		_, err = part.Write([]byte(testAudioData))
		require.NoError(t, err)
		writer.Close()

		// Create request with auth context
		req := httptest.NewRequest(http.MethodPost, "/reports", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))

		contentToBeStreamed := []inferenceService.ContentChanPayload{
			{Key: "_id", Value: testReportID},
			{Key: reports.Subjective, Value: "test subjective content"},
			{Key: reports.Objective, Value: "test objective content"},
			{Key: reports.AssessmentAndPlan, Value: "test assessment content"},
			{Key: reports.Summary, Value: "test summary content"},
		}

		// Setup mock expectations
		mockInference.On("GenerateReportPipeline",
			mock.MatchedBy(func(req *inferenceService.ReportRequest) bool {
				return req.ProviderID == testUserID &&
					string(req.AudioBytes) == testAudioData
			}),
			mock.AnythingOfType("chan inferenceService.ContentChanPayload"),
		).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan inferenceService.ContentChanPayload)

			for _, content := range contentToBeStreamed {
				ch <- content
			}

			close(ch)
		}).Return(nil).Once()

		rr := httptest.NewRecorder()
		handler.GenerateReport(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/x-ndjson", rr.Header().Get("Content-Type"))

		// Parse streamed responses
		var updates []inferenceService.ContentChanPayload
		decoder := json.NewDecoder(rr.Body)
		for decoder.More() {
			var update inferenceService.ContentChanPayload
			err := decoder.Decode(&update)
			require.NoError(t, err)
			updates = append(updates, update)
		}

		assert.Len(t, updates, len(contentToBeStreamed))
		assert.Equal(t, testReportID, updates[0].Value)
		assert.True(t, reflect.DeepEqual(contentToBeStreamed, updates))

		mockInference.AssertExpectations(t)
	})

	t.Run("should return internal server error when GenerateReportPipeline fails", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add audio file
		fileWriter, err := writer.CreateFormFile("audio", "test.wav")
		require.NoError(t, err)
		_, err = fileWriter.Write([]byte(testAudioData))
		require.NoError(t, err)

		// Add metadata
		metadata := map[string]interface{}{
			"patientName": "John Doe",
			"timestamp":   time.Now(),
			"duration":    60,
		}

		metadataBytes, err := json.Marshal(metadata)
		require.NoError(t, err)
		err = writer.WriteField("metadata", string(metadataBytes))
		require.NoError(t, err)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/reports", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))

		// Setup mock to return error
		mockInference.On("GenerateReportPipeline",
			mock.AnythingOfType("*inferenceService.ReportRequest"),
			mock.AnythingOfType("chan inferenceService.ContentChanPayload"),
		).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan inferenceService.ContentChanPayload)
			close(ch)
		}).Return(errors.New("pipeline error")).Once()

		rr := httptest.NewRecorder()
		handler.GenerateReport(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "should return 500 status code")
		assert.Contains(t, rr.Body.String(), "error generating report", "should return error message")

		mockInference.AssertExpectations(t)
	})

	t.Run("should return bad request when metadata is invalid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		err := writer.WriteField("metadata", "invalid json")
		require.NoError(t, err)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/reports", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))

		rr := httptest.NewRecorder()
		handler.GenerateReport(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid metadata")
	})

	t.Run("should return bad request when audio file isn't supplied", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		metadata := map[string]interface{}{
			"providerID": testUserID,
			"name":       "Test Report",
			"timestamp":  time.Now(),
			"duration":   30,
		}

		metadataBytes, err := json.Marshal(metadata)
		require.NoError(t, err)

		err = writer.WriteField("metadata", string(metadataBytes))
		require.NoError(t, err)

		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/reports", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))

		rr := httptest.NewRecorder()
		handler.GenerateReport(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "failed to get audio")
	})
}

func TestRegenerateReport(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)

	t.Run("should regenerate report when request is valid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Test Report",
		}

		req := inferenceService.ReportRequest{
			ID:         testReportID,
			ProviderID: testUserID,
			Updates:    bson.D{{Key: "subjective.data", Value: "updated subjective content"}},
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		// Create request with auth context
		httpReq := httptest.NewRequest(http.MethodPost, "/reports/regenerate", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		// Setup mock expectations
		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		contentToBeStreamed := []inferenceService.ContentChanPayload{
			{Key: "subjective", Value: "regenerated subjective content"},
			{Key: "objective", Value: "regenerated objective content"},
		}

		mockInference.On("RegenerateReport",
			mock.Anything,
			mock.AnythingOfType("chan inferenceService.ContentChanPayload"),
			mock.MatchedBy(func(r *inferenceService.ReportRequest) bool {
				return r.ID == testReportID && r.ProviderID == testUserID
			}),
		).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan inferenceService.ContentChanPayload)
			for _, content := range contentToBeStreamed {
				ch <- content
			}
			close(ch)
		}).Return(nil).Once()

		handler.RegenerateReport(rr, httpReq)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/x-ndjson", rr.Header().Get("Content-Type"))

		// Parse streamed responses
		var updates []inferenceService.ContentChanPayload
		decoder := json.NewDecoder(rr.Body)
		for decoder.More() {
			var update inferenceService.ContentChanPayload
			err := decoder.Decode(&update)
			require.NoError(t, err)
			updates = append(updates, update)
		}

		assert.Equal(t, len(contentToBeStreamed), len(updates))
		assert.True(t, reflect.DeepEqual(contentToBeStreamed, updates))

		MockReportsStore.AssertExpectations(t)
		mockInference.AssertExpectations(t)
	})

	t.Run("should return unauthorized when user not authenticated", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := httptest.NewRequest(http.MethodPost, "/reports/regenerate", nil)
		rr := httptest.NewRecorder()

		handler.RegenerateReport(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		invalidBody := []byte(`{"invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/reports/regenerate", bytes.NewBuffer(invalidBody))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.RegenerateReport(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid ReportRequest Format")
	})

	t.Run("should return error when report doesn't exist", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := inferenceService.ReportRequest{
			ID:         testReportID,
			ProviderID: testUserID,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/regenerate", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(reports.Report{}, errors.New("report not found")).Once()

		handler.RegenerateReport(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error regenerating report")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return error when provider ID doesn't match", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: "different-user-id",
			Name:       "Test Report",
		}

		req := inferenceService.ReportRequest{
			ID:         testReportID,
			ProviderID: testUserID,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/regenerate", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		handler.RegenerateReport(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error regenerating report")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return error when regeneration fails", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Test Report",
		}

		req := inferenceService.ReportRequest{
			ID:         testReportID,
			ProviderID: testUserID,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/regenerate", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		mockInference.On("RegenerateReport",
			mock.Anything,
			mock.AnythingOfType("chan inferenceService.ContentChanPayload"),
			mock.Anything,
		).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan inferenceService.ContentChanPayload)
			close(ch)
		}).Return(errors.New("regeneration error")).Once()

		handler.RegenerateReport(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error regenerating report")

		MockReportsStore.AssertExpectations(t)
		mockInference.AssertExpectations(t)
	})
}

func TestGetReport(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)

	t.Run("should return report when request is valid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		expectedReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Test Report",
			Subjective: reports.ReportContent{
				Data:    "test subjective content",
				Loading: false,
			},
		}

		req := GetReportRequest{
			ReportID: testReportID,
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodGet, "/reports/GetReport", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(expectedReport, nil).Once()

		handler.GetReport(rr, httpReq)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response reports.Report
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, expectedReport.ProviderID, response.ProviderID)
		assert.Equal(t, expectedReport.Name, response.Name)
		assert.Equal(t, expectedReport.Subjective, response.Subjective)
		assert.Equal(t, expectedReport.Objective, response.Objective)

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when user not authenticated", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := httptest.NewRequest(http.MethodGet, "/reports/GetReport", nil)
		rr := httptest.NewRecorder()

		handler.GetReport(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		invalidBody := []byte(`{"invalid json`)
		req := httptest.NewRequest(http.MethodGet, "/reports/GetReport", bytes.NewBuffer(invalidBody))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.GetReport(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")
	})

	t.Run("should return not found when report doesn't exist", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := GetReportRequest{
			ReportID: testReportID,
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodGet, "/reports/"+testReportID, bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(reports.Report{}, errors.New("report not found")).Once()

		handler.GetReport(rr, httpReq)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "error fetching report")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when provider ID doesn't match", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: "different-user-id",
			Name:       "Test Report",
		}

		req := GetReportRequest{
			ReportID: testReportID,
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodGet, "/reports/"+testReportID, bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		handler.GetReport(rr, httpReq)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")

		MockReportsStore.AssertExpectations(t)
	})
}

func TestLearnStyle(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)

	testContent := "test content"

	t.Run("should learn style when request is valid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Test Report",
		}

		req := LearnStyleRequest{
			ReportID:       testReportID,
			ContentSection: reports.Subjective,
			Previous:        testContent,
			Current: testContent,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/LearnStyle", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		mockInference.On("LearnStyle",
			mock.Anything,
			testUserID,
			reports.Subjective,
			testContent,
		).Return(nil).Once()

		handler.LearnStyle(rr, httpReq)

		assert.Equal(t, http.StatusOK, rr.Code)

		MockReportsStore.AssertExpectations(t)
		mockInference.AssertExpectations(t)
	})

	t.Run("should return unauthorized when user not authenticated", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := httptest.NewRequest(http.MethodPost, "/reports/LearnStyle", nil)
		rr := httptest.NewRecorder()

		handler.LearnStyle(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		invalidBody := []byte(`{"invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/reports/LearnStyle", bytes.NewBuffer(invalidBody))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.LearnStyle(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")
	})

	t.Run("should return error when report doesn't exist", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := LearnStyleRequest{
			ReportID:       testReportID,
			ContentSection: reports.Subjective,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/LearnStyle", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(reports.Report{}, errors.New("report not found")).Once()

		handler.LearnStyle(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error fetching report")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when provider ID doesn't match", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: "different-user-id",
			Name:       "Test Report",
		}

		req := LearnStyleRequest{
			ReportID:       testReportID,
			ContentSection: reports.Subjective,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/LearnStyle", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		handler.LearnStyle(rr, httpReq)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return error when learning style fails", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Test Report",
		}

		req := LearnStyleRequest{
			ReportID:       testReportID,
			ContentSection: reports.Subjective,
			Previous:        testContent,
			Current: testContent,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/LearnStyle", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		mockInference.On("LearnStyle",
			mock.Anything,
			testUserID,
			reports.Subjective,
			testContent,
		).Return(errors.New("learning style failed")).Once()

		handler.LearnStyle(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error learning style")

		MockReportsStore.AssertExpectations(t)
		mockInference.AssertExpectations(t)
	})
}

func TestGetTranscript(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)

	t.Run("should return transcript when request is valid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		expectedTranscript := "test transcript"

		req := GetReportRequest{
			ReportID: testReportID,
		}

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodGet, "/reports/GetTranscript", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("GetTranscription", mock.Anything, testReportID).
			Return(testUserID, expectedTranscript, nil).Once()

		handler.GetTranscript(rr, httpReq)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, expectedTranscript, response)
		

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when user not authenticated", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := httptest.NewRequest(http.MethodGet, "/reports/GetTranscript", nil)
		rr := httptest.NewRecorder()

		handler.GetTranscript(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		invalidBody := []byte(`{"invalid json`)
		req := httptest.NewRequest(http.MethodGet, "/reports/GetTranscript", bytes.NewBuffer(invalidBody))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.GetTranscript(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")
	})
}

func TestChangeReportName(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)

	t.Run("should change report name when request is valid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Old Report Name",
		}

		req := ChangeNameRequest{
			ReportID: testReportID,
			NewName:  "New Report Name",
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/ChangeName", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		MockReportsStore.On("UpdateReport", mock.Anything, testReportID, bson.D{
			{Key: "name", Value: "New Report Name"},
		}).Return(nil).Once()

		handler.ChangeReportName(rr, httpReq)

		assert.Equal(t, http.StatusOK, rr.Code)

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when user not authenticated", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := httptest.NewRequest(http.MethodPost, "/reports/ChangeName", nil)
		rr := httptest.NewRecorder()

		handler.ChangeReportName(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		invalidBody := []byte(`{"invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/reports/ChangeName", bytes.NewBuffer(invalidBody))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.ChangeReportName(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return error when report doesn't exist", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := ChangeNameRequest{
			ReportID: testReportID,
			NewName:  "New Report Name",
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/ChangeName", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(reports.Report{}, errors.New("error changing report name")).Once()

		handler.ChangeReportName(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error changing report name")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when provider ID doesn't match", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: "different-user-id",
			Name:       "Old Report Name",
		}

		req := ChangeNameRequest{
			ReportID: testReportID,
			NewName:  "New Report Name",
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/ChangeName", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		handler.ChangeReportName(rr, httpReq)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Old Report Name",
		}

		req := ChangeNameRequest{
			ReportID: testReportID,
			NewName:  "New Report Name",
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/ChangeName", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		MockReportsStore.On("UpdateReport", mock.Anything, testReportID, bson.D{
			{Key: "name", Value: "New Report Name"},
		}).Return(errors.New("update failed")).Once()

		handler.ChangeReportName(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error updating report")

		MockReportsStore.AssertExpectations(t)
	})
}

func TestUpdateContentData(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)

	t.Run("should update report when request is valid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Test Report",
		}

		req := UpdateContentData{
			ReportID:       testReportID,
			ContentSection: reports.Subjective,
			Content:        "updated subjective content",
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/UpdateReport", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		MockReportsStore.On("UpdateReport", mock.Anything, testReportID, bson.D{bson.E{Key: req.ContentSection, Value: bson.D{bson.E{Key: "data", Value: req.Content}}}}).Return(nil).Once()

		handler.UpdateContentSection(rr, httpReq)

		assert.Equal(t, http.StatusOK, rr.Code)

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when user not authenticated", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := httptest.NewRequest(http.MethodPost, "/reports/UpdateReport", nil)
		rr := httptest.NewRecorder()

		handler.UpdateContentSection(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		invalidBody := []byte(`{"invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/reports/UpdateReport", bytes.NewBuffer(invalidBody))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.UpdateContentSection(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return error when report doesn't exist", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := UpdateContentData{
			ReportID:       testReportID,
			ContentSection: reports.Subjective,
			Content:        "updated subjective content",
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(reports.Report{}, errors.New("report not found")).Once()

		httpReq := httptest.NewRequest(http.MethodPost, "/reports/UpdateReport", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.UpdateContentSection(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error updating report")

		MockReportsStore.AssertExpectations(t)
	})
}

func TestDeleteReport(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)

	t.Run("should delete report when request is valid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Test Report",
		}

		req := DeleteReportRequest{
			ReportIDs: []string{testReportID},
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodDelete, "/reports/DeleteReport", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		MockReportsStore.On("Delete", mock.Anything, testReportID).
			Return(nil).Once()

		handler.DeleteReport(rr, httpReq)

		assert.Equal(t, http.StatusOK, rr.Code)

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when user not authenticated", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := httptest.NewRequest(http.MethodDelete, "/reports/DeleteReport", nil)
		rr := httptest.NewRecorder()

		handler.DeleteReport(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		invalidBody := []byte(`{"invalid json`)
		req := httptest.NewRequest(http.MethodDelete, "/reports/DeleteReport", bytes.NewBuffer(invalidBody))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.DeleteReport(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")
	})

	t.Run("should return not found when report doesn't exist", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		req := DeleteReportRequest{
			ReportIDs: []string{testReportID},
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodDelete, "/reports/DeleteReport", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(reports.Report{}, errors.New("report not found")).Once()

		handler.DeleteReport(rr, httpReq)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized access to report")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when provider ID doesn't match", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: "different-user-id",
			Name:       "Test Report",
		}

		req := DeleteReportRequest{
			ReportIDs: []string{testReportID},
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodDelete, "/reports/DeleteReport", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		handler.DeleteReport(rr, httpReq)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")

		MockReportsStore.AssertExpectations(t)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		MockReportsStore := new(reports.MockReportsStore)
		mockInference := new(inferenceService.MockInferenceService)
		mockUser := new(user.MockUserStore)
		handler := NewReportsHandler(MockReportsStore, mockInference, mockUser, logger)

		existingReport := reports.Report{
			ProviderID: testUserID,
			Name:       "Test Report",
		}

		req := DeleteReportRequest{
			ReportIDs: []string{testReportID},
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodDelete, "/reports/DeleteReport", bytes.NewBuffer(body))
		httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		MockReportsStore.On("Get", mock.Anything, testReportID).
			Return(existingReport, nil).Once()

		MockReportsStore.On("Delete", mock.Anything, testReportID).
			Return(errors.New("delete failed")).Once()

		handler.DeleteReport(rr, httpReq)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error deleting report")

		MockReportsStore.AssertExpectations(t)
	})
}
