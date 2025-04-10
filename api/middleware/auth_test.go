package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var auth *AuthMiddleware

func TestMain(m *testing.M) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()
	auth = NewAuthMiddleware(os.Getenv("JWT_SECRET"), logger, "dev")
	code := m.Run()
	os.Exit(code)
}

type env struct {
	t                    *testing.T
	previousAccessToken  string
	previousRefreshToken string
	userID               string
	req                  *http.Request
}

func setUpEnv(t *testing.T) env {
	claims := &Claims{
		UserID: "userID",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	assert.NoError(t, err)

	refreshClaims := &Claims{
		UserID: "userID",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  AccessToken,
		Value: accessTokenString,
	})
	req.AddCookie(&http.Cookie{
		Name:  RefreshToken,
		Value: refreshTokenString,
	})

	return env{
		t:                    t,
		previousAccessToken:  accessTokenString,
		previousRefreshToken: refreshTokenString,
		userID:               "userID",
		req:                  req,
	}
}

func (e *env) setToken(tokenType string, expiresAt time.Time, value string) {
	if tokenType != AccessToken && tokenType != RefreshToken {
		e.t.Fatalf("invalid token name: %s", tokenType)
	}

	var tokenString string
	if value != "" {
		tokenString = value
	} else {
		claims := &Claims{
			UserID: "userID",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		genToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		assert.NoError(e.t, err)
		tokenString = genToken
	}

	e.req = httptest.NewRequest("GET", "/", nil)
	e.req.AddCookie(&http.Cookie{
		Name:  tokenType,
		Value: tokenString,
	})
	if tokenType == AccessToken {
		e.previousAccessToken = tokenString
		e.req.AddCookie(&http.Cookie{
			Name:  RefreshToken,
			Value: e.previousRefreshToken,
		})
	} else {
		e.previousRefreshToken = tokenString
		e.req.AddCookie(&http.Cookie{
			Name:  AccessToken,
			Value: e.previousAccessToken,
		})
	}
}

func getTokenFromCookies(t *testing.T, rr *httptest.ResponseRecorder, tokenType string) http.Cookie {
	if tokenType != AccessToken && tokenType != RefreshToken {
		t.Fatalf("invalid token name: %s", tokenType)
	}
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == tokenType {
			return *cookie
		}
	}
	t.Fatalf("cookie not found: %s", tokenType)
	return http.Cookie{}
}

func TestAuthMiddleware_ValidAccessToken(t *testing.T) {
	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetProviderIDFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "userID", userID)
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("should return 200 status code when access token is valid", func(t *testing.T) {
		env := setUpEnv(t)

		rr := httptest.NewRecorder()

		requestCookie, err := env.req.Cookie(AccessToken)
		assert.Nil(t, err)

		handler.ServeHTTP(rr, env.req)
		assert.Equal(t, http.StatusOK, rr.Code)

		newAccessTokenCookie := getTokenFromCookies(t, rr, AccessToken)

		assert.True(t, requestCookie.Value == newAccessTokenCookie.Value, "should not get a new access_token")
	})

	t.Run("should return 200 status code when access token is expired, and refresh token is valid", func(t *testing.T) {
		env := setUpEnv(t)
		env.setToken(AccessToken, time.Now().Add(-time.Hour), "")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, env.req)
		assert.Equal(t, http.StatusOK, rr.Code)

		newAccessTokenCookie := getTokenFromCookies(t, rr, AccessToken)
		assert.False(t, env.previousAccessToken == newAccessTokenCookie.Value, "should get new access_token")
	})
}

func TestAuthMiddleware_Failure(t *testing.T) {
	t.Run("should fail with invalid access token", func(t *testing.T) {
		env := setUpEnv(t)
		env.setToken(AccessToken, time.Now().Add(time.Hour), "invalid.token.signature")
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  AccessToken,
			Value: "invalid.token.signature",
		})
		req.AddCookie(&http.Cookie{
			Name:  RefreshToken,
			Value: "valid.refresh.token",
		})

		rr := httptest.NewRecorder()
		handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should fail with expired access token and invalid refresh token", func(t *testing.T) {
		env := setUpEnv(t)
		env.setToken(AccessToken, time.Now().Add(-time.Hour), "")
		env.setToken(RefreshToken, time.Now().Add(time.Hour), "invalid.refresh.token")

		rr := httptest.NewRecorder()
		handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		}))

		handler.ServeHTTP(rr, env.req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should fail when both tokens are expired", func(t *testing.T) {
		env := setUpEnv(t)
		env.setToken(AccessToken, time.Now().Add(-time.Hour), "")
		env.setToken(RefreshToken, time.Now().Add(-time.Hour), "")

		rr := httptest.NewRecorder()
		handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		}))

		handler.ServeHTTP(rr, env.req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestGenerateInitialTokens_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	err := auth.AttachInitialTokens(rr, "userID")
	assert.NoError(t, err)

	cookie := getTokenFromCookies(t, rr, AccessToken)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, "userID", claims.UserID)
}

func TestGetProviderIDFromContext(t *testing.T) {
	t.Run("valid user ID in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), CtxKeyUserID, "userID")
		userID, ok := GetProviderIDFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, "userID", userID)
	})
	t.Run("no user ID in context", func(t *testing.T) {
		ctx := context.Background()
		userID, ok := GetProviderIDFromContext(ctx)
		assert.False(t, ok)
		assert.Equal(t, "", userID)
	})

	t.Run("invalid type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), CtxKeyUserID, 1)
		userID, ok := GetProviderIDFromContext(ctx)
		assert.False(t, ok)
		assert.Equal(t, "", userID)
	})

}
