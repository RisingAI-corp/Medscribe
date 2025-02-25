package integrationtests

// import (
// 	userhandler "Medscribe/api/handlers/userHandler"
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"Medscribe/api/utils"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// const (
// 	testName     = "Test User"
// 	testEmail    = "test@example.com"
// 	testPassword = "password123"
// )

// func TestUserRoutesIntegration(t *testing.T) {
// 	t.Run("should create user and store in database when signup is successful", func(t *testing.T) {
// 		testEnv, err := SetupTestEnv()
// 		require.NoError(t, err)
// 		t.Cleanup(func() {
// 			err := testEnv.CleanupTestData()
// 			assert.NoError(t, err)
// 			err = testEnv.Disconnect()
// 			assert.NoError(t, err)
// 		})

// 		reqBody := userhandler.SignUpRequest{
// 			Name:     testName,
// 			Email:    testEmail,
// 			Password: testPassword,
// 		}
// 		body, err := json.Marshal(reqBody)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPost, "/user/signup", bytes.NewBuffer(body))
// 		req.Header.Set("Content-Type", "application/json")
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusCreated, rr.Code)

// 		var response userhandler.AuthResponse
// 		err = json.NewDecoder(rr.Body).Decode(&response)
// 		require.NoError(t, err)
// 		assert.NotEmpty(t, response.ID)

// 		user, err := testEnv.GetTestUser(response.ID)
// 		require.NoError(t, err)
// 		assert.Equal(t, testName, user.Name)
// 		assert.Equal(t, testEmail, user.Email)

// 		cookies := rr.Result().Cookies()
// 		_, err = utils.GetAccessToken(cookies)
// 		assert.Nil(t, err)
// 		_, err = utils.GetRefreshToken(cookies)
// 		assert.Nil(t, err)

// 	})

// 	t.Run("should return conflict when email already exists", func(t *testing.T) {
// 		testEnv, err := SetupTestEnv()
// 		require.NoError(t, err)
// 		t.Cleanup(func() {
// 			err := testEnv.CleanupTestData()
// 			assert.NoError(t, err)
// 			err = testEnv.Disconnect()
// 			assert.NoError(t, err)
// 		})

// 		userID, err := testEnv.CreateTestUser(testName, testEmail, testPassword)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, userID)

// 		reqBody := userhandler.SignUpRequest{
// 			Name:     "Another User",
// 			Email:    testEmail,
// 			Password: "different123",
// 		}
// 		body, err := json.Marshal(reqBody)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPost, "/user/signup", bytes.NewBuffer(body))
// 		req.Header.Set("Content-Type", "application/json")
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)
// 		assert.Equal(t, http.StatusConflict, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "email already in use")
// 	})

// 	t.Run("should authenticate user and return reports when login is successful", func(t *testing.T) {
// 		testEnv, err := SetupTestEnv()
// 		require.NoError(t, err)
// 		t.Cleanup(func() {
// 			err := testEnv.CleanupTestData()
// 			assert.NoError(t, err)
// 			err = testEnv.Disconnect()
// 			assert.NoError(t, err)
// 		})
// 		userID, err := testEnv.CreateTestUser(testName, testEmail, testPassword)
// 		require.NoError(t, err)

// 		reportID, err := testEnv.CreateTestReport(userID)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, reportID)

// 		reqBody := userhandler.LoginRequest{
// 			Email:    testEmail,
// 			Password: testPassword,
// 		}
// 		body, err := json.Marshal(reqBody)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(body))
// 		req.Header.Set("Content-Type", "application/json")
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)

// 		var response userhandler.AuthResponse
// 		err = json.NewDecoder(rr.Body).Decode(&response)
// 		require.NoError(t, err)
// 		assert.Equal(t, userID, response.ID)

// 		cookies := rr.Result().Cookies()
// 		_, err = utils.GetAccessToken(cookies)
// 		assert.Nil(t, err)
// 		_, err = utils.GetRefreshToken(cookies)
// 		assert.Nil(t, err)

// 		reports, err := testEnv.GetAllReports(userID)
// 		require.NoError(t, err)
// 		assert.Equal(t, len(reports), len(response.Reports))
// 		assert.Equal(t, response.Reports[0].ID, reports[0].ID)
// 	})

// 	t.Run("should return unauthorized when credentials are invalid", func(t *testing.T) {
// 		testEnv, err := SetupTestEnv()
// 		require.NoError(t, err)
// 		t.Cleanup(func() {
// 			err := testEnv.CleanupTestData()
// 			assert.NoError(t, err)
// 			err = testEnv.Disconnect()
// 			assert.NoError(t, err)
// 		})

// 		reqBody := userhandler.LoginRequest{
// 			Email:    "wrong@email.com",
// 			Password: "wrongpass",
// 		}
// 		body, err := json.Marshal(reqBody)
// 		require.NoError(t, err)

// 		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(body))
// 		req.Header.Set("Content-Type", "application/json")
// 		rr := httptest.NewRecorder()

// 		testEnv.Router.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusUnauthorized, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "invalid credentials")
// 	})
// }
