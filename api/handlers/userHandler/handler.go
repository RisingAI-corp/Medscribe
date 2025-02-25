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
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Email           string           `json:"email"`
	Reports         []reports.Report `json:"reports"`
	SubjectiveStyle string           `json:"subjectiveStyle"`
	ObjectiveStyle  string           `json:"objectiveStyle"`
	AssessmentStyle string           `json:"assessmentStyle"`
	PlanningStyle   string           `json:"planningStyle"`
	SummaryStyle    string           `json:"summaryStyle"`
}

type UserHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetMe(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userStore user.UserStore
	reports   reports.Reports
	logger    *zap.Logger
}

func NewUserHandler(userStore user.UserStore, reports reports.Reports, logger *zap.Logger) UserHandler {
	return &userHandler{
		userStore: userStore,
		reports:   reports,
		logger:    logger,
	}
}

func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ProviderID, err := h.userStore.Put(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == fmt.Sprintf("user already exists with this email: %s", req.Email) {
			http.Error(w, "email already in use", http.StatusConflict)
			return
		}
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	if err := middleware.GenerateInitialTokens(w, ProviderID); err != nil {
		http.Error(w, "failed to generate auth tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(AuthResponse{
		ID:    ProviderID,
		Name:  req.Name,
		Email: req.Email,
	}); err != nil {
		http.Error(w, "error encoding auth response", http.StatusInternalServerError)
		return
	}
	if err := middleware.AttachInitialTokens(w, ProviderID); err != nil {
		//TODO: log specifically this error
		http.Error(w, "error creating authentication", http.StatusInternalServerError)
	}
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := h.userStore.GetByAuth(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	userID := user.ID.Hex()

	if err := middleware.GenerateInitialTokens(w, userID); err != nil {
		http.Error(w, "failed to generate auth tokens", http.StatusInternalServerError)
		return
	}

	reports, err := h.reports.GetAll(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(AuthResponse{
		ID:              user.ID.Hex(),
		Name:            user.Name,
		Email:           user.Email,
		Reports:         reports,
		SubjectiveStyle: user.SubjectiveStyle,
		ObjectiveStyle:  user.ObjectiveStyle,
		AssessmentStyle: user.AssessmentStyle,
		PlanningStyle:   user.PlanningStyle,
		SummaryStyle:    user.SummaryStyle,
	}); err != nil {
		http.Error(w, "error encoding auth response", http.StatusInternalServerError)
		return
	}
	if err := middleware.AttachInitialTokens(w, userID); err != nil {
		//TODO: log specifically this error
		http.Error(w, "error creating authentication", http.StatusInternalServerError)
	}
}

func (h *userHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ProviderID, ok := middleware.GetProviderIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.userStore.Get(r.Context(), fmt.Sprint(ProviderID))
	if err != nil {
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}
	userID := user.ID.Hex()
	reports, err := h.reports.GetAll(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch reports", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(AuthResponse{
		ID:              user.ID.Hex(),
		Name:            user.Name,
		Email:           user.Email,
		Reports:         reports,
		SubjectiveStyle: user.SubjectiveStyle,
		ObjectiveStyle:  user.ObjectiveStyle,
		AssessmentStyle: user.AssessmentStyle,
		PlanningStyle:   user.PlanningStyle,
		SummaryStyle:    user.SummaryStyle,
	})
	if err != nil {
		http.Error(w, "failed to fetch reports", http.StatusInternalServerError)
		return
	}
}
