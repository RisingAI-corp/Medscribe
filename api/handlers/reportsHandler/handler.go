package reportsHandler

import (
	"Medscribe/api/middleware"
	inferenceService "Medscribe/inference/service"
	contextLogger "Medscribe/logger"
	"Medscribe/reports"
	"Medscribe/user"
	"context"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type ReportsHandler interface {
	GenerateReport(w http.ResponseWriter, r *http.Request)
	RegenerateReport(w http.ResponseWriter, r *http.Request)
	LearnStyle(w http.ResponseWriter, r *http.Request)
	ChangeReportName(w http.ResponseWriter, r *http.Request)
	UpdateContentSection(w http.ResponseWriter, r *http.Request)
	GetReport(w http.ResponseWriter, r *http.Request)
	DeleteReport(w http.ResponseWriter, r *http.Request)
	GetTranscript(w http.ResponseWriter, r *http.Request)
	MarkRead(w http.ResponseWriter, r *http.Request)
	MarkUnread(w http.ResponseWriter, r *http.Request)
}
type GetReportRequest struct {
	ReportID string `json:"reportID"`
}
type DeleteReportRequest struct {
	ReportIDs []string `json:"reportIDs"`
}

type LearnStyleRequest struct {
	ReportID       string `json:"reportID"`
	ContentSection string `json:"contentSection"`
	Previous       string `json:"previous"`
	Current        string `json:"current"`
}

type ChangeNameRequest struct {
	ReportID string `json:"reportID"`
	NewName  string `json:"newName"`
}

type UpdateContentData struct {
	ReportID       string `json:"reportID"`
	ContentSection string `json:"contentSection"`
	Content        string `json:"content"`
}

type reportsHandler struct {
	reportsService   reports.Reports
	inferenceService inferenceService.InferenceService
	userStore        user.UserStore
	logger           *zap.Logger
}

type ReadStatusRequest struct {
	ReportID string `json:"reportID"`
	Opened   bool   `json:"opened"`
}

func NewReportsHandler(reportsService reports.Reports, inferenceService inferenceService.InferenceService, userStore user.UserStore, logger *zap.Logger) ReportsHandler {
	return &reportsHandler{
		reportsService:   reportsService,
		inferenceService: inferenceService,
		userStore:        userStore,
		logger:           logger,
	}
}

func (h *reportsHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	// Use the logger from context
	logger := contextLogger.FromCtx(r.Context())

	// Authorization check
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		logger.Error("User is not authorized", zap.String("UserID", userID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Log that report generation is starting
	logger.Info("Starting report generation", zap.String("UserID", userID))

	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logger.Error("Failed to parse form", zap.Error(err))
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Decode the metadata into our request struct
	var req inferenceService.ReportRequest
	err := json.NewDecoder(strings.NewReader(r.FormValue("metadata"))).Decode(&req)
	if err != nil {
		logger.Error("Invalid metadata", zap.Error(err))
		http.Error(w, "invalid metadata", http.StatusBadRequest)
		return
	}
	req.ProviderID = userID

	// Get the audio file from the form
	file, _, err := r.FormFile("audio")
	if err != nil {
		logger.Error("Failed to get audio file", zap.Error(err))
		http.Error(w, "failed to get audio file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the audio bytes
	audioBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read audio", zap.Error(err))
		http.Error(w, "failed to read audio", http.StatusInternalServerError)
		return
	}
	req.AudioBytes = audioBytes

	// Set up SSE headers for streaming
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create channels for content and errors
	contentChan := make(chan inferenceService.ContentChanPayload)
	errChan := make(chan error)

	// Start the report generation pipeline in a goroutine
	go func() {
		// Log pipeline start
		logger.Info("Report generation pipeline started", zap.String("UserID", userID))
		if err := h.inferenceService.GenerateReportPipeline(r.Context(), &req, contentChan); err != nil {
			logger.Error("Error during report generation pipeline", zap.Error(err))
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Assert that the ResponseWriter supports streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Stream the content as it's received from the content channel
	for content := range contentChan {
		if err := json.NewEncoder(w).Encode(content); err != nil {
			logger.Error("Error writing response", zap.Error(err))
			http.Error(w, "error writing response", http.StatusInternalServerError)
			return
		}
		flusher.Flush()
	}

	// Check for errors from the pipeline
	if err := <-errChan; err != nil {
		logger.Error("Error generating report", zap.Error(err))
		http.Error(w, "error generating report", http.StatusInternalServerError)
		return
	}

	// Log successful report generation
	logger.Info("Report generation completed successfully", zap.String("UserID", userID))
}
func (h *reportsHandler) RegenerateReport(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	// Validate user authentication and log the start of regeneration.
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	logger.Info("Starting report regeneration", zap.String("UserID", userID))

	var req inferenceService.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid ReportRequest Format", zap.Error(err))
		http.Error(w, "invalid ReportRequest Format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req.ProviderID = userID

	// Verify that the report exists and the user is authorized to regenerate it.
	report, err := h.reportsService.Get(r.Context(), req.ID)
	if err != nil {
		logger.Error("Error regenerating report: failed to fetch report", zap.Error(err))
		http.Error(w, "error regenerating report", http.StatusInternalServerError)
		return
	}
	if report.ProviderID != userID {
		logger.Error("Unauthorized access to report", zap.String("UserID", userID), zap.String("ReportID", req.ID))
		http.Error(w, "error regenerating report", http.StatusInternalServerError)
		return
	}

	// Set up headers for streaming response.
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	contentChan := make(chan inferenceService.ContentChanPayload, 6)
	errChan := make(chan error)
	go func() {
		logger.Info("Regeneration pipeline started", zap.String("UserID", userID), zap.String("ReportID", req.ID))
		if err := h.inferenceService.RegenerateReport(r.Context(), contentChan, &req); err != nil {
			logger.Error("Error during report regeneration pipeline", zap.Error(err))
			errChan <- err
			return
		}
		errChan <- nil
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	for content := range contentChan {
		if err := json.NewEncoder(w).Encode(content); err != nil {
			logger.Error("Error writing response", zap.Error(err))
			http.Error(w, "error writing response", http.StatusInternalServerError)
			return
		}
		flusher.Flush()
	}

	if err := <-errChan; err != nil {
		logger.Error("Error regenerating report", zap.Error(err))
		http.Error(w, "error regenerating report", http.StatusInternalServerError)
		return
	}

	logger.Info("Report regeneration completed successfully", zap.String("UserID", userID), zap.String("ReportID", req.ID))
}

func (h *reportsHandler) LearnStyle(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	providerID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req LearnStyleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("Starting LearnStyle process", zap.String("ProviderID", providerID), zap.String("ReportID", req.ReportID))

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		logger.Error("Error fetching report", zap.Error(err))
		http.Error(w, "error fetching report", http.StatusInternalServerError)
		return
	}
	if report.ProviderID != providerID {
		logger.Error("Unauthorized access to report", zap.String("ProviderID", providerID), zap.String("ReportID", req.ReportID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err = h.inferenceService.LearnStyle(r.Context(), providerID, req.ContentSection, req.Previous, req.Current); err != nil {
		logger.Error("Learning style failed", zap.Error(err))
		http.Error(w, "error learning style", http.StatusInternalServerError)
		return
	}

	logger.Info("LearnStyle process completed successfully", zap.String("ProviderID", providerID), zap.String("ReportID", req.ReportID))
}

func (h *reportsHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req GetReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("Fetching report", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		logger.Error("Error fetching report", zap.Error(err))
		http.Error(w, "error fetching report", http.StatusNotFound)
		return
	}
	if report.ProviderID != userID {
		logger.Error("Unauthorized access to report", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(w).Encode(report); err != nil {
		logger.Error("Error encoding report", zap.Error(err))
		http.Error(w, "error encoding report", http.StatusInternalServerError)
		return
	}

	logger.Info("Report fetched successfully", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
}

func (h *reportsHandler) ChangeReportName(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		logger.Warn("Unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req ChangeNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("Attempting to change report name", zap.String("UserID", userID), zap.String("ReportID", req.ReportID), zap.String("NewName", req.NewName))

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		logger.Error("Error fetching report", zap.Error(err))
		http.Error(w, "error changing report name", http.StatusInternalServerError)
		return
	}
	if report.ProviderID != userID {
		logger.Error("Unauthorized access to report", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	updates := bson.D{bson.E{Key: reports.Name, Value: req.NewName}}
	if err = h.reportsService.UpdateReport(r.Context(), req.ReportID, updates); err != nil {
		logger.Error("Error updating report", zap.Error(err))
		http.Error(w, "error updating report", http.StatusInternalServerError)
		return
	}

	logger.Info("Report name changed successfully", zap.String("ReportID", req.ReportID), zap.String("NewName", req.NewName))
}

func (h *reportsHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		logger.Warn("Unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req ReadStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("Attempting to mark report as read", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		logger.Error("Error fetching report", zap.Error(err))
		http.Error(w, "error marking report as read", http.StatusInternalServerError)
		return
	}
	if report.ProviderID != userID {
		logger.Error("Unauthorized access to report", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err = h.reportsService.MarkRead(r.Context(), req.ReportID); err != nil {
		logger.Error("Error marking report as read", zap.Error(err))
		http.Error(w, "error marking report as read", http.StatusInternalServerError)
		return
	}

	logger.Info("Report marked as read successfully", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
}

func (h *reportsHandler) MarkUnread(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		logger.Warn("Unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req ReadStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("Attempting to mark report as unread", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		logger.Error("Error fetching report", zap.Error(err))
		http.Error(w, "error marking report as unread", http.StatusInternalServerError)
		return
	}
	if report.ProviderID != userID {
		logger.Error("Unauthorized access to report", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err = h.reportsService.MarkUnread(r.Context(), req.ReportID); err != nil {
		logger.Error("Error marking report as unread", zap.Error(err))
		http.Error(w, "error marking report as unread", http.StatusInternalServerError)
		return
	}

	logger.Info("Report marked as unread successfully", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
}

func (h *reportsHandler) GetTranscript(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		logger.Warn("Unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req GetReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("Attempting to get transcript", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))

	providerID, transcript, err := h.reportsService.GetTranscription(r.Context(), req.ReportID)
	if err != nil {
		logger.Error("Error fetching report transcript", zap.Error(err))
		http.Error(w, "error fetching report", http.StatusInternalServerError)
		return
	}
	if providerID != userID {
		logger.Error("Unauthorized access to report", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(w).Encode(transcript); err != nil {
		logger.Error("Error encoding transcript", zap.Error(err))
		http.Error(w, "error encoding transcript", http.StatusInternalServerError)
		return
	}

	logger.Info("Transcript fetched successfully", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
}

func (h *reportsHandler) UpdateContentSection(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		logger.Warn("Unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateContentData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("Attempting to update content section", zap.String("UserID", userID), zap.String("ReportID", req.ReportID), zap.String("ContentSection", req.ContentSection))

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		logger.Error("Error fetching report", zap.Error(err))
		http.Error(w, "error updating report", http.StatusInternalServerError)
		return
	}
	if report.ProviderID != userID {
		logger.Error("Unauthorized access to report", zap.String("UserID", userID), zap.String("ReportID", req.ReportID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	updates := bson.D{bson.E{Key: req.ContentSection, Value: bson.D{bson.E{Key: "data", Value: req.Content}}}}
	if err = h.reportsService.UpdateReport(r.Context(), req.ReportID, updates); err != nil {
		logger.Error("Error updating report", zap.Error(err))
		http.Error(w, "error updating report", http.StatusInternalServerError)
		return
	}

	logger.Info("Content section updated successfully", zap.String("ReportID", req.ReportID), zap.String("ContentSection", req.ContentSection))
}

func (h *reportsHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())

	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		logger.Warn("Unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req DeleteReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("Attempting to delete reports", zap.String("UserID", userID), zap.Strings("ReportIDs", req.ReportIDs))

	if err := h.verifyReportBelongsToProvider(r.Context(), userID, req.ReportIDs...); err != nil {
		logger.Error("Unauthorized access to report during deletion", zap.Error(err))
		http.Error(w, "unauthorized access to report", http.StatusUnauthorized)
		return
	}

	if err := h.runDeleteReports(r.Context(), req.ReportIDs...); err != nil {
		logger.Error("Error deleting report", zap.Error(err))
		http.Error(w, "error deleting report", http.StatusInternalServerError)
		return
	}

	logger.Info("Reports deleted successfully", zap.Strings("ReportIDs", req.ReportIDs))
	w.WriteHeader(http.StatusOK)
}

func (h *reportsHandler) verifyReportBelongsToProvider(r context.Context, providerID string, reportIDs ...string) error {
	logger := contextLogger.FromCtx(r)
	for _, reportID := range reportIDs {
		report, err := h.reportsService.Get(r, reportID)
		if err != nil {
			logger.Error("Error verifying report ownership", zap.String("ReportID", reportID), zap.Error(err))
			return fmt.Errorf("error verifying report belongs to provider: %s", providerID)
		}
		if providerID != report.ProviderID {
			logger.Error("Unauthorized access detected", zap.String("ReportID", reportID), zap.String("ProviderID", providerID))
			return fmt.Errorf("unauthorized access to report: %s", reportID)
		}
	}
	return nil
}

func (h *reportsHandler) runDeleteReports(r context.Context, reportIDs ...string) error {
	logger := contextLogger.FromCtx(r)
	for _, reportID := range reportIDs {
		if err := h.reportsService.Delete(r, reportID); err != nil {
			logger.Error("Error deleting report", zap.String("ReportID", reportID), zap.Error(err))
			return err
		}
	}
	return nil
}
