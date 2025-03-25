package reportsHandler

import (
	"Medscribe/api/middleware"
	inferenceService "Medscribe/inference/service"
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

func NewReportsHandler(reportsService reports.Reports, inferenceService inferenceService.InferenceService, userStore user.UserStore, logger *zap.Logger) ReportsHandler {
	return &reportsHandler{
		reportsService:   reportsService,
		inferenceService: inferenceService,
		userStore:        userStore,
		logger:           logger,
	}
}

func (h *reportsHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		h.logger.Error("user is not authorized: ", zap.String("UserID: ", userID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Error("failed to parse form", zap.Error(err))
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req inferenceService.ReportRequest
	err := json.NewDecoder(strings.NewReader(r.FormValue("metadata"))).Decode(&req)
	if err != nil {
		h.logger.Error("invalid metadata", zap.Error(err))
		http.Error(w, "invalid metadata", http.StatusBadRequest)
		return
	}

	req.ProviderID = userID

	file, _, err := r.FormFile("audio")
	if err != nil {
		h.logger.Error("failed to get audio file", zap.Error(err))
		http.Error(w, "failed to get audio file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	audioBytes, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error("failed to read audio", zap.Error(err))
		http.Error(w, "failed to read audio", http.StatusInternalServerError)
		return
	}
	req.AudioBytes = audioBytes

	// Set up SSE headers for streaming
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	contentChan := make(chan inferenceService.ContentChanPayload)
	errChan := make(chan error)
	go func() {
		if err := h.inferenceService.GenerateReportPipeline(context.Background(), &req, contentChan); err != nil {
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
			h.logger.Error("error writing response", zap.Error(err))
			http.Error(w, "error writing response", http.StatusInternalServerError)
		}
		flusher.Flush()
	}
	if err := <-errChan; err != nil {
		h.logger.Error("error generating report", zap.Error(err))
		http.Error(w, "error generating report", http.StatusInternalServerError)
	}
}

func (h *reportsHandler) RegenerateReport(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req inferenceService.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid ReportRequest Format", zap.Error(err))
		http.Error(w, "invalid ReportRequest Format", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	req.ProviderID = userID

	report, err := h.reportsService.Get(r.Context(), req.ID)
	if err != nil {
		h.logger.Error("error regenerating report", zap.Error(err))
		http.Error(w, "error regenerating report ", http.StatusInternalServerError)
		return
	}
	if report.ProviderID != userID {
		h.logger.Error("unauthorized access to report", zap.Error(err))
		http.Error(w, "error regenerating report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	contentChan := make(chan inferenceService.ContentChanPayload, 6)

	errChan := make(chan error)
	go func() {
		if err := h.inferenceService.RegenerateReport(context.Background(), contentChan, &req); err != nil {
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
			h.logger.Error("error writing response", zap.Error(err))
			http.Error(w, "error writing response", http.StatusInternalServerError)
		}

		flusher.Flush()
	}

	if err := <-errChan; err != nil {
		h.logger.Error("error regenerating report", zap.Error(err))
		http.Error(w, "error regenerating report", http.StatusInternalServerError)
	}
}

func (h *reportsHandler) LearnStyle(w http.ResponseWriter, r *http.Request) {
	providerID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req LearnStyleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		h.logger.Error("error fetching report", zap.Error(err))
		http.Error(w, "error fetching report", http.StatusInternalServerError)
		return
	}
	if report.ProviderID != providerID {
		h.logger.Error("unauthorized access to report", zap.Error(err))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.inferenceService.LearnStyle(r.Context(), providerID, req.ContentSection, req.Previous, req.Current)
	if err != nil {
		h.logger.Error("learning style failed", zap.Error(err))
		http.Error(w, "error learning style", http.StatusInternalServerError)
		return
	}
}

func (h *reportsHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req GetReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		h.logger.Error("error fetching report", zap.Error(err))
		http.Error(w, "error fetching report", http.StatusNotFound)
		return
	}
	if report.ProviderID != userID {
		h.logger.Error("unauthorized access to report", zap.Error(err))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(w).Encode(report); err != nil {
		h.logger.Error("error encoding report", zap.Error(err))
		http.Error(w, "error encoding report", http.StatusInternalServerError)
		return
	}
}

func (h *reportsHandler) ChangeReportName(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		h.logger.Warn("unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req ChangeNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.logger.Info("attempting to change report name", zap.String("userID", userID), zap.String("reportID", req.ReportID), zap.String("newName", req.NewName))

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		h.logger.Error("error fetching report", zap.Error(err))
		http.Error(w, "error changing report name", http.StatusInternalServerError)
		return
	}

	if report.ProviderID != userID {
		h.logger.Error("unauthorized access to report", zap.Error(err))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	updates := bson.D{bson.E{Key: reports.Name, Value: req.NewName}}

	err = h.reportsService.UpdateReport(r.Context(), req.ReportID, updates)
	if err != nil {
		h.logger.Error("error updating report", zap.Error(err))
		http.Error(w, "error updating report", http.StatusInternalServerError)
		return
	}

	h.logger.Info("report name changed successfully", zap.String("reportID", req.ReportID), zap.String("newName", req.NewName))
}

func (h *reportsHandler) GetTranscript(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		h.logger.Warn("unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req GetReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.logger.Info("attempting to get transcript", zap.String("userID", userID), zap.String("reportID", req.ReportID))

	providerID, transcript, err := h.reportsService.GetTranscription(r.Context(), req.ReportID)
	if err != nil {
		h.logger.Error("error fetching report", zap.Error(err))
		http.Error(w, "error fetching report", http.StatusInternalServerError)
		return
	}

	if providerID != userID {
		h.logger.Error("unauthorized access to report", zap.Error(err))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := json.NewEncoder(w).Encode(transcript); err != nil {
		h.logger.Error("error encoding transcript", zap.Error(err))
		http.Error(w, "error encoding transcript", http.StatusInternalServerError)
		return
	}
}

func (h *reportsHandler) UpdateContentSection(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		h.logger.Warn("unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateContentData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.logger.Info("attempting to update content section", zap.String("userID", userID), zap.String("reportID", req.ReportID), zap.String("contentSection", req.ContentSection))

	report, err := h.reportsService.Get(r.Context(), req.ReportID)
	if err != nil {
		h.logger.Error("error fetching report", zap.Error(err))
		http.Error(w, "error updating report", http.StatusInternalServerError)
		return
	}

	if report.ProviderID != userID {
		h.logger.Error("unauthorized access to report", zap.Error(err))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	updates := bson.D{bson.E{Key: req.ContentSection, Value: bson.D{bson.E{Key: "data", Value: req.Content}}}}

	err = h.reportsService.UpdateReport(r.Context(), req.ReportID, updates)
	if err != nil {
		h.logger.Error("error updating report", zap.Error(err))
		http.Error(w, "error updating report", http.StatusInternalServerError)
		return
	}

	h.logger.Info("content section updated successfully", zap.String("reportID", req.ReportID), zap.String("contentSection", req.ContentSection))
}

func (h *reportsHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		h.logger.Warn("unauthorized access attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req DeleteReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.logger.Info("attempting to delete reports", zap.String("userID", userID), zap.Strings("reportIDs", req.ReportIDs))

	if err := h.verifyReportBelongsToProvider(r.Context(), userID, req.ReportIDs...); err != nil {
		h.logger.Error("unauthorized access to report", zap.Error(err))
		http.Error(w, "unauthorized access to report", http.StatusUnauthorized)
		return
	}

	if err := h.runDeleteReports(r.Context(), req.ReportIDs...); err != nil {
		h.logger.Error("error deleting report", zap.Error(err))
		http.Error(w, "error deleting report", http.StatusInternalServerError)
		return
	}

	h.logger.Info("reports deleted successfully", zap.Strings("reportIDs", req.ReportIDs))
	w.WriteHeader(http.StatusOK)
}

func (h *reportsHandler) verifyReportBelongsToProvider(r context.Context, providerID string, reportIDs ...string) error {
	for _, reportID := range reportIDs {
		report, err := h.reportsService.Get(r, reportID)
		if err != nil {
			h.logger.Error("error deleting report", zap.Error(err))
			return fmt.Errorf("error verifying report belongs to provider: %s", providerID)
		}
		if providerID != report.ProviderID {
			return fmt.Errorf("unauthorized access to report: %s", reportID)
		}
	}

	return nil
}

func (h *reportsHandler) runDeleteReports(r context.Context, reportIDs ...string) error {
	for _, reportID := range reportIDs {
		err := h.reportsService.Delete(r, reportID)
		if err != nil {
			h.logger.Error("error deleting report", zap.Error(err))
			return err
		}
	}

	return nil
}
