package userhandler

import (
	"Medscribe/api/middleware"
	"Medscribe/reports"
	"Medscribe/user"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Request/Response types
type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required, password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	ID                       string           `json:"id"`
	Name                     string           `json:"name"`
	Email                    string           `json:"email"`
	Reports                  []reports.Report `json:"reports"`
	SubjectiveStyle          string           `json:"subjectiveStyle"`
	ObjectiveStyle           string           `json:"objectiveStyle"`
	AssessmentAndPlanStyle   string           `json:"assessmentAndPlanStyle"`
	PatientInstructionsStyle string           `json:"patientInstructionsStyle"`
	PlanningStyle            string           `json:"planningStyle"`
	SummaryStyle             string           `json:"summaryStyle"`
	UserID                   string           `json:"userID"` // Added UserID field
}

type UserHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetMe(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userStore      user.UserStore
	reports        reports.Reports
	logger         *zap.Logger
	authMiddleware middleware.AuthMiddleware
}

func NewUserHandler(userStore user.UserStore, reports reports.Reports, logger *zap.Logger, authMiddleware middleware.AuthMiddleware) UserHandler {
	return &userHandler{
		userStore:      userStore,
		reports:        reports,
		logger:         logger,
		authMiddleware: authMiddleware,
	}
}

func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Error decoding signup request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ProviderID, err := h.userStore.Put(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		h.logger.Error("Error creating user", zap.Error(err))
		if err.Error() == fmt.Sprintf("user already exists with this email: %s", req.Email) {
			http.Error(w, "email already in use", http.StatusConflict)
			return
		}
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	if err := h.authMiddleware.AttachInitialTokens(w, ProviderID); err != nil {
		h.logger.Error("Error generating auth tokens", zap.Error(err))
		http.Error(w, "failed to generate auth tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(AuthResponse{
		ID:      ProviderID,
		Name:    req.Name,
		Email:   req.Email,
		UserID:  ProviderID, // Added UserID here
		Reports: []reports.Report{},
	}); err != nil {
		h.logger.Error("Error encoding auth response", zap.Error(err))
		http.Error(w, "error encoding auth response", http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Error decoding login request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := h.userStore.GetByAuth(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Error authenticating user", zap.Error(err))
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	userID := user.ID.Hex()

	if err := h.authMiddleware.AttachInitialTokens(w, userID); err != nil {
		h.logger.Error("Error generating auth tokens", zap.Error(err))
		http.Error(w, "failed to generate auth tokens", http.StatusInternalServerError)
		return
	}

	reports, err := h.reports.GetAll(r.Context(), userID)
	if err != nil {
		h.logger.Error("Error fetching reports", zap.Error(err))
		http.Error(w, "failed to get reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(AuthResponse{
		ID:                       user.ID.Hex(),
		Name:                     user.Name,
		Email:                    user.Email,
		Reports:                  reports,
		SubjectiveStyle:          user.SubjectiveStyle,
		ObjectiveStyle:           user.ObjectiveStyle,
		AssessmentAndPlanStyle:   user.AssessmentAndPlanStyle,
		PatientInstructionsStyle: user.PatientInstructionsStyle,
		SummaryStyle:             user.SummaryStyle,
		UserID:                   userID, // Added UserID here
	}); err != nil {
		h.logger.Error("Error encoding auth response", zap.Error(err))
		http.Error(w, "error encoding auth response", http.StatusInternalServerError)
		return
	}
	if err := h.authMiddleware.AttachInitialTokens(w, userID); err != nil {
		h.logger.Error("Error creating authentication tokens", zap.Error(err)) // Log specific error
		http.Error(w, "error creating authentication", http.StatusInternalServerError)
	}
}

func (h *userHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ProviderID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		h.logger.Error("Unauthorized request: ProviderID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.userStore.Get(r.Context(), fmt.Sprint(ProviderID))
	if err != nil {
		h.logger.Error("Error fetching user", zap.Error(err))
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}
	userID := user.ID.Hex()
	reports, err := h.reports.GetAll(r.Context(), userID)
	if err != nil {
		h.logger.Error("Error fetching reports", zap.Error(err))
		http.Error(w, "failed to fetch reports", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(AuthResponse{
		ID:                       user.ID.Hex(),
		Name:                     user.Name,
		Email:                    user.Email,
		Reports:                  reports,
		SubjectiveStyle:          user.SubjectiveStyle,
		ObjectiveStyle:           user.ObjectiveStyle,
		AssessmentAndPlanStyle:   user.AssessmentAndPlanStyle,
		PatientInstructionsStyle: user.PatientInstructionsStyle,
		SummaryStyle:             user.SummaryStyle,
		UserID:                   userID, // Added UserID here
	})
	if err != nil {
		h.logger.Error("Error encoding getMe response", zap.Error(err))
		http.Error(w, "failed to fetch reports", http.StatusInternalServerError)
		return
	}
}
