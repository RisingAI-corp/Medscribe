package userhandler

import (
	"Medscribe/api/middleware"
	emailsender "Medscribe/emailService"
	contextLogger "Medscribe/logger"
	"Medscribe/reports"
	"Medscribe/user"
	verificationStore "Medscribe/verificationTokenStore"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type InitializingSignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required, password"`
}

type FinalizeSignUpRequest struct {
	Token  string `json:"token" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type InitializeForgotPasswordRequest struct {
	Token string `json:"token" validate:"required"`
	Email string `json:"email" validate:"required,email"`
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
type UpdateProfileSettingsRequest struct {
	Name                     string `json:"name" validate:"required"`
	CurrentPassword          string `json:"currentPassword" validate:"required"`
	NewPassword              string `json:"newPassword" validate:"required"`
}

type UserHandler interface {
	InitializeSighUp(w http.ResponseWriter, r *http.Request)
	FinalizeSignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetMe(w http.ResponseWriter, r *http.Request)
	UpdateProfileSettings(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userStore      user.UserStore
	reports        reports.Reports
	authMiddleware middleware.AuthMiddleware
	verificationStore verificationStore.VerificationStore
	emailSender 	emailsender.EmailSender
}

func NewUserHandler(userStore user.UserStore, reports reports.Reports, authMiddleware middleware.AuthMiddleware, verificationStore verificationStore.VerificationStore ,emailSender emailsender.EmailSender) UserHandler {
	return &userHandler{
		userStore:      userStore,
		reports:        reports,
		authMiddleware: authMiddleware,
		verificationStore: verificationStore,
		emailSender: 	emailSender,
	}
}

func (h *userHandler) InitializeSighUp(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())
	logger.Info("initializing signup process")
	w.Header().Set("Content-Type", "application/json")

	var req InitializingSignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode signup initialization request", zap.Error(err))
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	doesEmailExist, err := h.userStore.CheckEmailExistence(r.Context(), req.Email)
	if err != nil {
		logger.Error("failed to check email existence", zap.Error(err))
		http.Error(w, user.EmailAlreadyExistsError.Error(), http.StatusConflict)
		return
	}
	if doesEmailExist{
		logger.Warn("email already exists", zap.String("email", req.Email))
		http.Error(w, "email already exists", http.StatusConflict)
		return
	}

	token, err := verificationStore.GenerateRandomToken(6);
	if err != nil {
		logger.Error("failed to store verification document", zap.Error(err))
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	err = h.verificationStore.PutBufferedUserDocument(r.Context(),token ,req.Name, req.Email, req.Password)
	if err != nil {
		if err == verificationStore.ErrVerificationDocumentAlreadyExists {
			logger.Warn("Cannot generate Verification token this soon", zap.String("token", token))
			fmt.Println("BRONCO")
			http.Error(w, "Too Soon to generate Verification Token", http.StatusConflict)
			return
		}
		logger.Error("failed to store verification document", zap.Error(err))
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	err = h.emailSender.SendEmail(req.Email, "Please Verify your email!", fmt.Sprintf("Your verification token is: %s", token), emailsender.GenerateVerificationHTMLBody(token, "Verify your email"))
	if err != nil {
		logger.Error("failed to send verification email", zap.Error(err))
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}



func (h *userHandler) FinalizeSignUp(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())
	var req FinalizeSignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode signup request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	fmt.Println("Received token:", req.Token)

	BufferedProviderDocument, err := h.verificationStore.GetBufferedUserDocument(r.Context(), req.Token)
	if err != nil {
		if err == verificationStore.ErrVerificationDocumentNotFound {
			logger.Warn("verification token not found", zap.String("token", req.Token))
			http.Error(w, "verification token not found", http.StatusBadRequest)
			return
		}

		logger.Error("failed to get buffered user document", zap.Error(err))
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	err = h.verificationStore.Delete(r.Context(), req.Token)
	if err != nil {	
		logger.Error("failed to delete buffered user document", zap.Error(err))
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	// Log the start of putting the user into the database.
	logger.Info("starting to store new user in the database",
		zap.String("name", BufferedProviderDocument.Name),
		zap.String("email", BufferedProviderDocument.Email),
	)

	ProviderID, err := h.userStore.Put(r.Context(), BufferedProviderDocument.Name, BufferedProviderDocument.Email, BufferedProviderDocument.Password)
	if err != nil {
		if err.Error() == fmt.Sprintf("user already exists with this email: %s", BufferedProviderDocument.Email) {
			logger.Warn("user already exists",
				zap.String("email", BufferedProviderDocument.Email),
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
		ID:      ProviderID,
		Name:    BufferedProviderDocument.Name,
		Email:   BufferedProviderDocument.Email,
		UserID:  ProviderID,
		Reports: []reports.Report{},
	}); err != nil {
		logger.Error("failed to encode auth response", zap.Error(err))
		http.Error(w, "error encoding auth response", http.StatusInternalServerError)
		return
	}

	// Log a successful signup and response.
	logger.Info("user signup successful",
		zap.String("provider_id", ProviderID),
		zap.String("name", BufferedProviderDocument.Name),
		zap.String("email", BufferedProviderDocument.Email),
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

func (h *userHandler) InitializeForgotPassword(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())
	var req InitializeForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode forgot password request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()	

	err := h.verificationStore.PutResetPasswordDetails(r.Context(), req.Token, req.Email)
	if err != nil {	
		if err == verificationStore.ErrVerificationDocumentAlreadyExists {
			logger.Warn("Cannot generate Verification token this soon", zap.String("token", req.Token))
			return 
		}
		logger.Error("failed to store reset password details", zap.Error(err))
		http.Error(w, "error Putting reset password details", http.StatusInternalServerError)
		return
	}
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", r.Host, req.Token)
	logger.Info("sending reset password email",
		zap.String("email", req.Email),
		zap.String("reset_link", resetLink),
	)
	err = h.emailSender.SendEmail(req.Email, "Please Reset your password!", fmt.Sprintf("Your reset password token is: %s", req.Token), emailsender.GeneratePasswordResetHTMLBody(resetLink))
	if err != nil {
		logger.Error("failed to send reset password email", zap.Error(err))
		http.Error(w, "error sending reset password email", http.StatusInternalServerError)
		return	
	}
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

func (h *userHandler) UpdateProfileSettings(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())
	ProviderID, ok := middleware.GetProviderIDFromContext(r.Context())
	userIDStr := fmt.Sprint(ProviderID)
	if !ok {
		logger.Warn("unauthorized access attempt to update user profile settings")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req UpdateProfileSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode update profile settings request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	logger.Info("updating user profile settings",
		zap.String("user_id", userIDStr),

	)
	err := h.userStore.UpdateProfileSettings(r.Context(), userIDStr, req.Name, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "invalid current password" {
			logger.Warn("invalid current password provided",
				zap.String("user_id", userIDStr),
			)
			http.Error(w, "invalid current password", http.StatusPreconditionFailed)
			return
		}
		logger.Error("failed to update user profile settings", zap.Error(err))
	http.Error(w, "failed to update user profile settings", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}


func (h *userHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := contextLogger.FromCtx(r.Context())
	logger.Info("handling user logout")

	// Clear the access token cookie
	accessTokenCookie := &http.Cookie{
		Name:     middleware.AccessToken,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, accessTokenCookie)

	// Clear the refresh token cookie
	refreshTokenCookie := &http.Cookie{
		Name:     middleware.RefreshToken,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, refreshTokenCookie)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	logger.Info("user logout successful")
}