package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type contextKey string

const UserIDKey contextKey = "userID"

const AccessToken = "access_token"
const RefreshToken = "refresh_token"

// Default token durations
const (
	DefaultAccessTokenDuration  = 15 * time.Minute
	DefaultRefreshTokenDuration = 7 * 24 * time.Hour // 1 week
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims *Claims
		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		//attempt to get refresh token
		claims, err = verifyToken(refreshCookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		regenerateAccessToken := false
		var accessToken string

		//check if access token doesn't exist and if it does verify it
		accessCookie, err := r.Cookie(AccessToken)
		if err != nil {
			regenerateAccessToken = true
		} else {
			claims, err = verifyToken(accessCookie.Value)
			if err != nil {
				regenerateAccessToken = true
			}
			accessToken = accessCookie.Value
		}

		if regenerateAccessToken {
			accessToken, err = GenerateAccessToken(claims.UserID)
			if err != nil {
				//TODO: make sure to add reason why
				http.Error(w, "error occurred generating access token", http.StatusInternalServerError)
			}
		}

		setCookies(w, accessToken, refreshCookie.Value)

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims
		expirationTime, err := claims.GetExpirationTime()
		if err != nil {
			return "", fmt.Errorf("error getting expiration time: %w", err)
		}

		if expirationTime.Before(jwt.NewNumericDate(time.Now()).Time) {
			return "", fmt.Errorf("token has expired")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return claims, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return claims, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GenerateAccessToken(userID string) (string, error) {
	accessToken, err := generateToken(userID, DefaultAccessTokenDuration)
	if err != nil {
		return "", fmt.Errorf("error generating access token: %w", err)
	}

	return accessToken, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	refreshToken, err := generateToken(userID, DefaultRefreshTokenDuration)
	if err != nil {
		return "", fmt.Errorf("error generating refresh token: %w", err)
	}

	return refreshToken, nil
}

func AttachInitialTokens(w http.ResponseWriter, userID string) error {
	accessToken, err := generateToken(userID, DefaultAccessTokenDuration)
	if err != nil {
		return fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := generateToken(userID, DefaultRefreshTokenDuration)
	if err != nil {
		return fmt.Errorf("error generating refresh token: %w", err)
	}
	setCookies(w, accessToken, refreshToken)

	return nil
}

func generateToken(userID string, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

func setCookies(w http.ResponseWriter, access_token, refresh_token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessToken,
		Value:    access_token,
		Path:     "/",
		MaxAge:   int(DefaultAccessTokenDuration.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     RefreshToken,
		Value:    refresh_token,
		Path:     "/",
		MaxAge:   int(DefaultRefreshTokenDuration.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

}

// GenerateInitialTokens generates the first set of tokens for a user (e.g., after login)
func GenerateInitialTokens(w http.ResponseWriter, userID string) error {
	accessToken, err := GenerateAccessToken(userID)
	if err != nil {
		return err
	}
	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return err
	}

	setCookies(w, accessToken, refreshToken)
	return nil
}

// GetProviderIDFromContext retrieves the user ID from the context
func GetProviderIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", false
	}
	return userID, true
}
