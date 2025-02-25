package userhandler

import (
	"Medscribe/api/middleware"
	"Medscribe/reports"
	"Medscribe/user"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// Test constants
const (
	testName     = "Test User"
	testEmail    = "test@example.com"
	testPassword = "password123"
	testUserID   = "507f1f77bcf86cd799439011"
)

func checkCookie(cookies []*http.Cookie, cookieName string) (string, bool) {
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie.Value, true
		}
	}
	return "", false
}

func TestSignUp(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)
	defer func() {
		err := logger.Sync()
		if err != nil {
			t.Errorf("failed to sync logger: %v", err)
		}
	}()
	t.Run("should create user when credentials are valid", func(t *testing.T) {
		mockStore := new(user.MockUserStore)
		mockReportsStore := new(reports.MockReports)
		handler := NewUserHandler(mockStore, mockReportsStore, logger)

		mockStore.On("Put", mock.Anything, testName, testEmail, testPassword).
			Return(testUserID, nil)

		reqBody := SignUpRequest{
			Name:     testName,
			Email:    testEmail,
			Password: testPassword,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.SignUp(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response AuthResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, testUserID, response.ID)
		assert.Equal(t, testName, response.Name)
		assert.Equal(t, testEmail, response.Email)
		assert.Empty(t, response.Reports)

		accessToken, ok := checkCookie(rr.Result().Cookies(), "access_token")
		assert.NotEmpty(t, accessToken)
		assert.True(t, ok)

		refreshToken, ok := checkCookie(rr.Result().Cookies(), "refresh_token")
		assert.NotEmpty(t, refreshToken)
		assert.True(t, ok)

		mockStore.AssertExpectations(t)
	})

	t.Run("should return conflict when email already exists", func(t *testing.T) {
		mockStore := new(user.MockUserStore)
		mockReportsStore := new(reports.MockReports)
		handler := NewUserHandler(mockStore, mockReportsStore, logger)

		mockStore.On("Put", mock.Anything, testName, testEmail, testPassword).
			Return("", fmt.Errorf("user already exists with this email: %s", testEmail))

		reqBody := SignUpRequest{
			Name:     testName,
			Email:    testEmail,
			Password: testPassword,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.SignUp(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
		assert.Contains(t, rr.Body.String(), "email already in use")

		mockStore.AssertExpectations(t)
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		mockStore := new(user.MockUserStore)
		mockReportsStore := new(reports.MockReports)
		handler := NewUserHandler(mockStore, mockReportsStore, logger)

		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		handler.SignUp(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")
	})
}

func TestLogin(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)
	defer func() {
		err := logger.Sync()
		if err != nil {
			t.Errorf("failed to sync logger: %v", err)
		}
	}()
	t.Run("should authenticate user when credentials are valid", func(t *testing.T) {
		mockStore := new(user.MockUserStore)
		mockReportsStore := new(reports.MockReports)
		handler := NewUserHandler(mockStore, mockReportsStore, logger)

		objectID, err := primitive.ObjectIDFromHex(testUserID)
		require.NoError(t, err)

		mockUser := user.User{
			ID:    objectID,
			Name:  testName,
			Email: testEmail,
		}

		mockStore.On("GetByAuth", mock.Anything, testEmail, testPassword).
			Return(mockUser, nil)

		mockReports := []reports.Report{
			{
				ID:         objectID,
				Name:       "Test Report",
				ProviderID: testUserID,
			},
		}

		mockReportsStore.On("GetAll", mock.Anything, testUserID).
			Return(mockReports, nil).
			Once()

		reqBody := LoginRequest{
			Email:    testEmail,
			Password: testPassword,
		}
		body, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response AuthResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response contents
		assert.Equal(t, testUserID, response.ID)
		assert.Equal(t, testName, response.Name)
		assert.Equal(t, testEmail, response.Email)
		assert.Len(t, response.Reports, 1)
		assert.Equal(t, mockReports[0].ID.Hex(), response.Reports[0].ID.Hex())
		assert.Equal(t, mockReports[0].Name, response.Reports[0].Name)
		assert.Equal(t, mockReports[0].ProviderID, response.Reports[0].ProviderID)

		accessToken, ok := checkCookie(rr.Result().Cookies(), "access_token")
		assert.NotEmpty(t, accessToken)
		assert.True(t, ok)

		refreshToken, ok := checkCookie(rr.Result().Cookies(), "refresh_token")
		assert.NotEmpty(t, refreshToken)
		assert.True(t, ok)

		mockStore.AssertExpectations(t)
		mockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when credentials are invalid", func(t *testing.T) {
		mockStore := new(user.MockUserStore)
		mockReportsStore := new(reports.MockReports)
		handler := NewUserHandler(mockStore, mockReportsStore, logger)

		mockStore.On("GetByAuth", mock.Anything, testEmail, testPassword).
			Return(user.User{}, fmt.Errorf("incorrect authentication credentials"))

		reqBody := LoginRequest{
			Email:    testEmail,
			Password: testPassword,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid credentials")

		mockStore.AssertExpectations(t)
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		mockStore := new(user.MockUserStore)
		mockReportsStore := new(reports.MockReports)
		handler := NewUserHandler(mockStore, mockReportsStore, logger)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request body")
	})
}

func TestGetMe(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)
	defer func() {
		err := logger.Sync()
		if err != nil {
			t.Errorf("failed to sync logger: %v", err)
		}
	}()

	t.Run("should return user data when authenticated", func(t *testing.T) {
		mockStore := new(user.MockUserStore)
		mockReportsStore := new(reports.MockReports)
		handler := NewUserHandler(mockStore, mockReportsStore, logger)

		objectID, err := primitive.ObjectIDFromHex(testUserID)
		if err != nil {
			t.Fatalf("failed to convert testUserID to ObjectID: %v", err)
		}
		mockUser := user.User{
			ID:    objectID,
			Name:  testName,
			Email: testEmail,
			Styles: user.Styles{
				SubjectiveStyle: "Test Subjective Style", // Add test style values
				ObjectiveStyle:  "Test Objective Style",
				AssessmentStyle: "Test Assessment Style",
				PlanningStyle:   "Test Planning Style",
				SummaryStyle:    "Test Summary Style",
			},
		}

		// Update mock expectations with correct arguments
		mockStore.On("Get", mock.Anything, testUserID).
			Return(mockUser, nil).
			Once()

		mockReports := []reports.Report{
			{
				ID:         objectID,
				Name:       "Test Report",
				ProviderID: testUserID,
			},
		}

		mockReportsStore.On("GetAll", mock.Anything, testUserID).
			Return(mockReports, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, testUserID))
		rr := httptest.NewRecorder()

		handler.GetMe(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response AuthResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response contents, including styles
		assert.Equal(t, testUserID, response.ID)
		assert.Len(t, response.Reports, 1)
		assert.Equal(t, mockReports[0].ID.Hex(), response.Reports[0].ID.Hex())
		assert.Equal(t, mockReports[0].Name, response.Reports[0].Name)
		assert.Equal(t, mockReports[0].ProviderID, response.Reports[0].ProviderID)

		assert.Equal(t, "Test Subjective Style", response.SubjectiveStyle) // Assert styles
		assert.Equal(t, "Test Objective Style", response.ObjectiveStyle)
		assert.Equal(t, "Test Assessment Style", response.AssessmentStyle)
		assert.Equal(t, "Test Planning Style", response.PlanningStyle)
		assert.Equal(t, "Test Summary Style", response.SummaryStyle)

		mockStore.AssertExpectations(t)
		mockReportsStore.AssertExpectations(t)
	})

	t.Run("should return unauthorized when user ID not in context", func(t *testing.T) {
		mockStore := new(user.MockUserStore)
		mockReportsStore := new(reports.MockReports)
		handler := NewUserHandler(mockStore, mockReportsStore, logger)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.GetMe(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})
}
