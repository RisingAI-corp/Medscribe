package integrationtests

import (
	"Medscribe/api/handlers/reportsHandler"
	"Medscribe/reports"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testUserName = "testuser"
const testUserEmail = "testuser@example.com"
const testUserPassword = "password"
const testPatientName = "testPatient"

func loadAudioFile(t *testing.T, filePath string) []byte {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	audioData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	return audioData
}

// func TestReportGenerationRoutesSuccess(t *testing.T) {
// 	testEnv, err := SetupTestEnv()
// 	if err != nil {
// 		t.Fatalf("Failed to setup test environment: %v", err)
// 	}
// 	t.Cleanup(func() {
// 		err := testEnv.CleanupTestData()
// 		assert.NoError(t, err)
// 		err = testEnv.Disconnect()
// 		assert.NoError(t, err)
// 	})

// 	userID, err := testEnv.CreateTestUser(testUserName, testUserEmail, testUserPassword)
// 	assert.NoError(t, err)

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	// Add metadata
// 	timestamp := time.Now()
// 	metadata := map[string]interface{}{
// 		"providerID":  userID,
// 		"patientName": testPatientName,
// 		"timestamp":   timestamp,
// 		"duration":    30,
// 	}

// 	metadataBytes, err := json.Marshal(metadata)
// 	assert.NoError(t, err)

// 	err = writer.WriteField("metadata", string(metadataBytes))
// 	require.NoError(t, err)

// 	// Add audio file
// 	part, err := writer.CreateFormFile("audio", "test.wav")
// 	require.NoError(t, err)

// 	audioFile := loadAudioFile(t, "../../testdata/sample1.wav")
// 	_, err = part.Write([]byte(audioFile))
// 	require.NoError(t, err)
// 	writer.Close()

// 	req := httptest.NewRequest(http.MethodPost, "/report/generate", body)
// 	req, err = testEnv.GenerateJWT(req, userID)

// 	require.NoError(t, err)

// 	req.Header.Set("Content-Type", writer.FormDataContentType())
// 	rr := httptest.NewRecorder()

// 	testEnv.Router.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	assert.Equal(t, "application/x-ndjson", rr.Header().Get("Content-Type"))

// 	// Parse streamed responses
// 	fields := []string{}
// 	reportID := ""
// 	hasEmptyValue := false
// 	decoder := json.NewDecoder(rr.Body)
// 	for decoder.More() {
// 		var update inferenceService.ContentChanPayload
// 		err := decoder.Decode(&update)
// 		require.NoError(t, err)
// 		fields = append(fields, update.Key)
// 		if update.Value == "" {
// 			hasEmptyValue = true
// 		}
// 		if update.Key == "_id" {
// 			reportID = update.Value.(string)
// 		}
// 	}

// 	expectedContentTypes := []string{"_id", reports.Subjective, reports.Objective, reports.AssessmentAndPlan, reports.PatientInstructions, reports.CondensedSummary, reports.SessionSummary, reports.Summary, reports.FinishedGenerating}

// 	assert.ElementsMatch(t, expectedContentTypes, fields, "expected content types do not match the fields")
// 	assert.True(t, !hasEmptyValue)

// 	time.Sleep(time.Second * 2)
// 	report, err := testEnv.GetTestReport(reportID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, testPatientName, report.Name)
// 	assert.Equal(t, userID, report.ProviderID)
// 	assert.Equal(t, float64(30), report.Duration)

// 	equalWithinMargin := func(t1, t2 time.Time, margin time.Duration) bool {
// 		return t1.Sub(t2) <= margin && t2.Sub(t1) <= margin
// 	}

// 	assert.True(t, equalWithinMargin(timestamp, report.TimeStamp.Time(), time.Second*1))

// 	assert.NotEmpty(t, report.Subjective.Data, "Subjective data should not be empty")
// 	assert.NotEmpty(t, report.Objective.Data, "Objective data should not be empty")
// 	assert.NotEmpty(t, report.AssessmentAndPlan.Data, "AssessmentAndPlan data should not be empty")
// 	assert.NotEmpty(t, report.PatientInstructions.Data, "Patient instructions data should not be empty")
// 	assert.NotEmpty(t, report.CondensedSummary, "Condensed summary data should not be empty")
// 	assert.NotEmpty(t, report.SessionSummary, "Session summary data should not be empty")
// 	assert.NotEmpty(t, report.Summary.Data, "Summary data should not be empty")
// 	assert.True(t, report.FinishedGenerating, "Report should still be in the process of generating")
// }

// func TestReportGenerationRoutesFailures(t *testing.T) {
// 	testEnv, err := SetupTestEnv()
// 	require.NoError(t, err)
// 	t.Cleanup(func() {
// 		err := testEnv.CleanupTestData()
// 		assert.NoError(t, err)
// 		err = testEnv.Disconnect()
// 		assert.NoError(t, err)
// 	})

// 	userID, err := testEnv.CreateTestUser(testUserName, testUserEmail, testUserPassword)
// 	require.NoError(t, err)

// 	t.Run("should return bad request when metadata field values are invalid", func(t *testing.T) {
// 		reqBody := &bytes.Buffer{}
// 		writer := multipart.NewWriter(reqBody)

// 		// Add invalid metadata
// 		metadata := map[string]interface{}{
// 			"providerID":  "",
// 			"patientName": testPatientName,
// 			"timestamp":   "invalid timestamp",
// 			"duration":    30,
// 			"reportContents": []inferenceService.ReportContentSection{
// 				{
// 					ContentType: reports.Subjective,
// 					Content:     "test subjective content",
// 				},
// 			},
// 		}

// 		metadataBytes, err := json.Marshal(metadata)
// 		require.NoError(t, err)

// 		err = writer.WriteField("metadata", string(metadataBytes))
// 		require.NoError(t, err)

// 		// Add a valid audio file
// 		part, err := writer.CreateFormFile("audio", "test.wav")
// 		require.NoError(t, err)

// 		audioFile := loadAudioFile(t, "../../testdata/sample1.wav")
// 		_, err = part.Write(audioFile)
// 		require.NoError(t, err)
// 		writer.Close()

// 		req := httptest.NewRequest(http.MethodPost, "/report/generate", reqBody)
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		req.Header.Set("Content-Type", writer.FormDataContentType())
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusBadRequest, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "invalid metadata")
// 	})

// 	t.Run("should return bad request when no audio file is supplied", func(t *testing.T) {
// 		reqBody := &bytes.Buffer{}
// 		writer := multipart.NewWriter(reqBody)

// 		metadata := map[string]interface{}{
// 			"providerID":  userID,
// 			"patientName": testPatientName,
// 			"timestamp":   time.Now(),
// 			"duration":    30,
// 			"reportContents": []inferenceService.ReportContentSection{
// 				{
// 					ContentType: reports.Subjective,
// 					Content:     "test subjective content",
// 				},
// 			},
// 		}

// 		metadataBytes, err := json.Marshal(metadata)
// 		require.NoError(t, err)

// 		err = writer.WriteField("metadata", string(metadataBytes))
// 		require.NoError(t, err)

// 		writer.Close()

// 		req := httptest.NewRequest(http.MethodPost, "/report/generate", reqBody)
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		req.Header.Set("Content-Type", writer.FormDataContentType())
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusBadRequest, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "failed to get audio file")
// 	})

// }

// func TestRegenerateReport(t *testing.T) {
// 	testEnv, err := SetupTestEnv()
// 	require.NoError(t, err)
// 	t.Cleanup(func() {
// 		err := testEnv.CleanupTestData()
// 		assert.NoError(t, err)
// 		err = testEnv.Disconnect()
// 		assert.NoError(t, err)
// 	})

// 	userID, err := testEnv.CreateTestUser(testUserName, testUserEmail, testUserPassword)
// 	require.NoError(t, err)

// 	// First, create a report to regenerate
// 	reportID, err := testEnv.CreateTestReport(userID)
// 	assert.NoError(t, err)

// 	t.Run("should successfully regenerate report", func(t *testing.T) {
// 		body := map[string]interface{}{
// 			"ID":                        reportID,
// 			"PatientName":               testPatientName,
// 			"AudioBytes":                []byte{},
// 			"TranscribedAudio":          "test transcribed audio",
// 			"ProviderID":                userID,
// 			"Timestamp":                 time.Now(),
// 			"Duration":                  30,
// 			"Updates":                   bson.D{{Key: "pronouns", Value: "HE"}},
// 			"SubjectiveContent":         "test subjective content",
// 			"ObjectiveContent":          "test objective content",
// 			"AssessmentAndPlanContent":  "test assessment and plan content",
// 			"PatientInstructionContent": "test patient instruction content",
// 			"SummaryContent":            "test summary content",
// 			"SubjectiveStyle":           "test subjective style",
// 			"ObjectiveStyle":            "test objective style",
// 			"AssessmentAndPlanStyle":    "test assessment and plan style",
// 			"SummaryStyle":              "test summary style",
// 			"SessionSummary":            "test session summary",
// 			"CondensedSummary":  "test condensed medical summary",
// 			"PatientInstructionsStyle": "test patient instructions style",
// 		}

// 		bodyBytes, err := json.Marshal(body)
// 		assert.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPatch, "/report/regenerate", bytes.NewBuffer(bodyBytes))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		assert.NoError(t, err)

// 		req.Header.Set("Content-Type", "application/json")
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusOK, rr.Code)

// 		// Process streamed responses
// 		fields := []string{}
// 		hasEmptyValue := false
// 		decoder := json.NewDecoder(rr.Body)
// 		for decoder.More() {
// 			var update inferenceService.ContentChanPayload
// 			err := decoder.Decode(&update)
// 			require.NoError(t, err)
// 			fields = append(fields, update.Key)
// 			if update.Value == "" {
// 				hasEmptyValue = true
// 			}
// 			if update.Key == "_id" {
// 				reportID = update.Value.(string)
// 			}
// 		}

// 		expectedContentTypes := []string{reports.Subjective, reports.Objective, reports.AssessmentAndPlan, reports.PatientInstructions, reports.CondensedSummary, reports.SessionSummary, reports.Summary, reports.FinishedGenerating}
// 		assert.ElementsMatch(t, expectedContentTypes, fields, "expected content types do not match the fields")
// 		assert.True(t, !hasEmptyValue)

// 		// Check that the report was updated correctly
// 		report, err := testEnv.GetTestReport(reportID)
// 		require.NoError(t, err)
// 		assert.NotEmpty(t, report.Subjective.Data)
// 	})

// 	t.Run("should return bad request when metadata field values are invalid", func(t *testing.T) {
// 		// Add invalid metadata
// 		metadata := map[string]interface{}{
// 			"providerID":  1, // invalid providerID type
// 			"patientName": testPatientName,
// 			"timestamp":   time.Now(),
// 			"duration":    30,
// 		}

// 		reportRequest, err := json.Marshal(metadata)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPatch, "/report/regenerate", bytes.NewBuffer(reportRequest))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		req.Header.Set("Content-Type", "application/json")
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusBadRequest, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "invalid ReportRequest Format")
// 	})

// 	t.Run("should return bad request when content types are invalid", func(t *testing.T) {
// 		// Add invalid metadata
// 		reportRequest := map[string]interface{}{
// 			"providerID":  userID,
// 			"patientName": testPatientName,
// 			"timestamp":   time.Now(),
// 			"duration":    30,
// 		}

// 		body, err := json.Marshal(reportRequest)
// 		assert.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPatch, "/report/regenerate", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		req.Header.Set("Content-Type", "application/json")
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusInternalServerError, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "error regenerating report")
// 	})

// // }

// func TestGetReport(t *testing.T) {
// 	testEnv, err := SetupTestEnv()
// 	require.NoError(t, err)
// 	t.Cleanup(func() {
// 		err := testEnv.CleanupTestData()
// 		assert.NoError(t, err)
// 		err = testEnv.Disconnect()
// 		assert.NoError(t, err)
// 	})

// 	userID, err := testEnv.CreateTestUser(testUserName, testUserEmail, testUserPassword)
// 	require.NoError(t, err)

// 	reportID, err := testEnv.CreateTestReport(userID)
// 	assert.NoError(t, err)

// 	t.Run("should successfully get report", func(t *testing.T) {
// 		getReportsRequest := reportsHandler.GetReportRequest{
// 			ReportID: reportID,
// 		}

// 		body, err := json.Marshal(getReportsRequest)
// 		assert.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPost, "/report/get", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusOK, rr.Code)

// 		var report inferenceService.ReportRequest
// 		err = json.Unmarshal(rr.Body.Bytes(), &report)

// 		assert.NoError(t, err)
// 		assert.Equal(t, reportID, report.ID)
// 	})

// 	t.Run("should return error when Bad Request", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodPost, "/report/get", nil)
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusBadRequest, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "invalid request body")

// 	})

// 	t.Run("should return not found when report does not exist", func(t *testing.T) {
// 		getReportsRequest := reportsHandler.GetReportRequest{
// 			ReportID: primitive.NewObjectID().Hex(),
// 		}

// 		body, err := json.Marshal(getReportsRequest)
// 		assert.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPost, "/report/get", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusNotFound, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "error fetching report")
// 	})

// 	t.Run("should return unauthorized when user is not authenticated", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodPost, "/report/get", nil)
// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusUnauthorized, rr.Code)
// 	})
// }

// func TestDeleteReport(t *testing.T) {
// 	testEnv, err := SetupTestEnv()
// 	require.NoError(t, err)
// 	t.Cleanup(func() {
// 		err := testEnv.CleanupTestData()
// 		assert.NoError(t, err)
// 		err = testEnv.Disconnect()
// 		assert.NoError(t, err)
// 	})

// 	userID, err := testEnv.CreateTestUser(testUserName, testUserEmail, testUserPassword)
// 	require.NoError(t, err)

// 	reportID, err := testEnv.CreateTestReport(userID)
// 	require.NoError(t, err)

// 	t.Run("should successfully delete report", func(t *testing.T) {
// 		deleteRequest := reportsHandler.DeleteReportRequest{
// 			ReportIDs: []string{reportID},
// 		}

// 		body, err := json.Marshal(deleteRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodDelete, "/report/delete", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)

// 		// Verify that the report has been deleted
// 		_, err = testEnv.GetTestReport(reportID)
// 		assert.Error(t, err)
// 	})

// 	t.Run("should return not found when report does not exist", func(t *testing.T) {
// 		deleteRequest := reportsHandler.DeleteReportRequest{
// 			ReportIDs: []string{reportID},
// 		}

// 		body, err := json.Marshal(deleteRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodDelete, "/report/delete", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusUnauthorized, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "unauthorized")
// 	})

// 	t.Run("should return unauthorized when user is not authenticated", func(t *testing.T) {
// 		deleteRequest := reportsHandler.DeleteReportRequest{
// 			ReportIDs: []string{reportID},
// 		}

// 		body, err := json.Marshal(deleteRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodDelete, "/report/delete", bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusUnauthorized, rr.Code)
// 	})
// }

// func TestGetTranscript(t *testing.T) {
// 	testEnv, err := SetupTestEnv()
// 	require.NoError(t, err)
// 	t.Cleanup(func() {
// 		err := testEnv.CleanupTestData()
// 		assert.NoError(t, err)
// 		err = testEnv.Disconnect()
// 		assert.NoError(t, err)
// 	})

// 	userID, err := testEnv.CreateTestUser(testUserName, testUserEmail, testUserPassword)
// 	require.NoError(t, err)

// 	reportID, err := testEnv.CreateTestReport(userID)
// 	require.NoError(t, err)

// 	t.Run("should return transcript when request is valid", func(t *testing.T) {
// 		getTranscriptRequest := reportsHandler.GetReportRequest{
// 			ReportID: reportID,
// 		}

// 		body, err := json.Marshal(getTranscriptRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPost, "/report/getTranscript", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)

// 		var response string
// 		err = json.NewDecoder(rr.Body).Decode(&response)
// 		fmt.Println(err, "eeee")
// 		assert.NoError(t, err)

// 	})

// 	t.Run("should return unauthorized when user is not authenticated", func(t *testing.T) {
// 		getTranscriptRequest := reportsHandler.GetReportRequest{
// 			ReportID: reportID,
// 		}

// 		body, err := json.Marshal(getTranscriptRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPost, "/report/getTranscript", bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusUnauthorized, rr.Code)
// 	})
// }

// func TestUpdateContentSection(t *testing.T) {
// 	testEnv, err := SetupTestEnv()
// 	require.NoError(t, err)
// 	t.Cleanup(func() {
// 		err := testEnv.CleanupTestData()
// 		assert.NoError(t, err)
// 		err = testEnv.Disconnect()
// 		assert.NoError(t, err)
// 	})

// 	userID, err := testEnv.CreateTestUser(testUserName, testUserEmail, testUserPassword)
// 	require.NoError(t, err)

// 	// First, create a report to update
// 	reportID, err := testEnv.CreateTestReport(userID)
// 	require.NoError(t, err)

// 	t.Run("should successfully update report", func(t *testing.T) {
// 		updateRequest := reportsHandler.UpdateContentData{
// 			ReportID:       reportID,
// 			ContentSection: reports.Subjective,
// 			Content:        "updated subjective content",
// 		}

// 		body, err := json.Marshal(updateRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPatch, "/report/updateContentSection", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)

// 		// Verify that the report has been updated
// 		report, err := testEnv.GetTestReport(reportID)
// 		require.NoError(t, err)
// 		assert.Equal(t, report.Subjective.Data, "updated subjective content")
// 	})

// 	t.Run("should return bad request when updates body is invalid", func(t *testing.T) {
// 		updateRequest := reportsHandler.UpdateContentData{
// 			ReportID:       reportID,
// 			ContentSection: "invalidSection",
// 			Content:        "updated subjective content",
// 		}

// 		body, err := json.Marshal(updateRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPatch, "/report/updateContentSection", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusInternalServerError, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "error updating report")
// 	})

// 	t.Run("should return unauthorized when user is not authenticated", func(t *testing.T) {
// 		updateRequest := reportsHandler.UpdateContentData{
// 			ReportID:       reportID,
// 			ContentSection: reports.Subjective,
// 			Content:        "updated subjective content",
// 		}

// 		body, err := json.Marshal(updateRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPatch, "/report/updateContentSection", bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusUnauthorized, rr.Code)
// 	})

// 	t.Run("should return internal server error when update fails", func(t *testing.T) {
// 		updateRequest := reportsHandler.UpdateContentData{
// 			ReportID:       reportID,
// 			ContentSection: "Invalid section",
// 			Content:        "updated subjective content",
// 		}

// 		body, err := json.Marshal(updateRequest)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPatch, "/report/updateContentSection", bytes.NewBuffer(body))
// 		req, err = testEnv.GenerateJWT(req, userID)
// 		require.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusInternalServerError, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "error updating report")
// 	})
// }

func TestLearnStyle(t *testing.T) {
	testEnv, err := SetupTestEnv()
	require.NoError(t, err)
	t.Cleanup(func() {
		err := testEnv.CleanupTestData()
		assert.NoError(t, err)
		err = testEnv.Disconnect()
		assert.NoError(t, err)
	})

	userID, err := testEnv.CreateTestUser(testUserName, testUserEmail, testUserPassword)
	require.NoError(t, err)

	// First, create a report to learn style from
	reportID, err := testEnv.CreateTestReport(userID)
	require.NoError(t, err)

	testContent := "test content"

	t.Run("should successfully learn style", func(t *testing.T) {
		learnStyleRequest := reportsHandler.LearnStyleRequest{
			ReportID:       reportID,
			ContentSection: reports.Subjective,
			Current:        testContent, // doesn't accept empty strings
			Previous:       "previous content",
		}

		body, err := json.Marshal(learnStyleRequest)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/report/learn-style", bytes.NewBuffer(body))
		req, err = testEnv.GenerateJWT(req, userID)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		testEnv.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/report/learn-style", nil)
		req, err = testEnv.GenerateJWT(req, userID)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		testEnv.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")
	})

	t.Run("should return unauthorized when user is not authenticated", func(t *testing.T) {
		learnStyleRequest := reportsHandler.LearnStyleRequest{
			ReportID:       reportID,
			ContentSection: reports.Subjective,
			Current:        "",
			Previous:       testContent,
		}

		body, err := json.Marshal(learnStyleRequest)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/report/learn-style", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		testEnv.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return internal server error when learning style fails", func(t *testing.T) {
		// Simulate a failure in the learning style operation
		learnStyleRequest := reportsHandler.LearnStyleRequest{
			ReportID:       reportID,
			ContentSection: reports.Subjective,
			Current:        "",
			Previous:       testContent,
		}

		body, err := json.Marshal(learnStyleRequest)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/report/learn-style", bytes.NewBuffer(body))
		req, err = testEnv.GenerateJWT(req, userID)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		testEnv.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error learning style")
	})
}
