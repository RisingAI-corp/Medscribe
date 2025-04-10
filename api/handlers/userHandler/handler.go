package userhandler

import (
	"Medscribe/api/middleware"
	contextLogger "Medscribe/logger"
	"Medscribe/reports"
	"Medscribe/user"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

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
	UserID                   string           `json:"userID"`
}

type UserHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetMe(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userStore      user.UserStore
	reports        reports.Reports
	authMiddleware middleware.AuthMiddleware
}

func NewUserHandler(userStore user.UserStore, reports reports.Reports, authMiddleware middleware.AuthMiddleware) UserHandler {
	return &userHandler{
		userStore:      userStore,
		reports:        reports,
		authMiddleware: authMiddleware,
	}
}
func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode signup request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the start of putting the user into the database.
	logger.Info("starting to store new user in the database",
		zap.String("name", req.Name),
		zap.String("email", req.Email),
	)

	ProviderID, err := h.userStore.Put(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == fmt.Sprintf("user already exists with this email: %s", req.Email) {
			logger.Warn("user already exists",
				zap.String("email", req.Email),
			)
			http.Error(w, "email already in use", http.StatusConflict)
			return
		}
		logger.Error("failed to store user in the database", zap.Error(err))
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	// Log a successful insertion into the database.
	logger.Info("user stored successfully in the database",
		zap.String("provider_id", ProviderID),
	)

	if err := h.authMiddleware.AttachInitialTokens(w, ProviderID); err != nil {
		logger.Error("failed to attach authentication tokens", zap.Error(err))
		http.Error(w, "failed to generate auth tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(AuthResponse{
		ID:          ProviderID,
		Name:        req.Name,
		Email:       req.Email,
		UserID:      ProviderID,
		Reports:     []reports.Report{},
	}); err != nil {
		logger.Error("failed to encode auth response", zap.Error(err))
		http.Error(w, "error encoding auth response", http.StatusInternalServerError)
		return
	}

	// Log a successful signup and response.
	logger.Info("user signup successful",
		zap.String("provider_id", ProviderID),
		zap.String("name", req.Name),
		zap.String("email", req.Email),
	)
}
func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode login request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info("attempting user login",
		zap.String("email", req.Email),
	)

	user, err := h.userStore.GetByAuth(r.Context(), req.Email, req.Password)
	if err != nil {
		logger.Warn("authentication failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return
	}
	userID := user.ID.Hex()

	logger.Info("user authenticated successfully",
		zap.String("provider_id", userID),
		zap.String("email", req.Email),
		zap.String("name", user.Name),
	)

	if err := h.authMiddleware.AttachInitialTokens(w, userID); err != nil {
		logger.Error("failed to attach authentication tokens", zap.Error(err))
		http.Error(w, "failed to generate auth tokens", http.StatusInternalServerError)
		return
	}

	reports, err := h.reports.GetAll(r.Context(), userID)
	if err != nil {
		logger.Error("failed to get reports for user", zap.Error(err), zap.String("user_id", userID))
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
		UserID:                   userID,
	}); err != nil {
		logger.Error("failed to encode auth response", zap.Error(err))
		http.Error(w, "error encoding auth response", http.StatusInternalServerError)
		return
	}

	// Log a successful login and response.
	logger.Info("user login successful",
		zap.String("user_id", userID),
		zap.String("email", user.Email),
		zap.String("name", user.Name),
	)

}

func (h *userHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())
	ProviderID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		logger.Warn("unauthorized access attempt to get user details")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userIDStr := fmt.Sprint(ProviderID)
	logger.Info("fetching user details", zap.String("user_id", userIDStr))

	user, err := h.userStore.Get(r.Context(), userIDStr)
	if err != nil {
		logger.Error("failed to fetch user details", zap.Error(err), zap.String("user_id", userIDStr))
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}
	userID := user.ID.Hex()
	reports, err := h.reports.GetAll(r.Context(), userID)
	if err != nil {
		logger.Error("failed to fetch reports for user", zap.Error(err), zap.String("user_id", userID))
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
		UserID:                   userID,
	})
	if err != nil {
		logger.Error("failed to encode user details response", zap.Error(err), zap.String("user_id", userID))
		http.Error(w, "failed to fetch reports", http.StatusInternalServerError)
		return
	}

	// Log a successful retrieval of user details.
	logger.Info("successfully retrieved user details",
		zap.String("user_id", userID),
		zap.String("name", user.Name),
		zap.String("email", user.Email),
		zap.Int("report_count", len(reports)),
	)
}