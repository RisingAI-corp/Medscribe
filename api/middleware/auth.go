package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)


const (
	UserIDKey    		 = "userID"
	AccessToken             = "access_token"
	RefreshToken            = "refresh_token"
)

const (
	DefaultAccessTokenDuration  = 15 * time.Minute
	DefaultRefreshTokenDuration = 7 * 24 * time.Hour
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	jwtSecret string
	logger    *zap.Logger
	env      string
	secure   bool


}

func NewAuthMiddleware(secret string, logger *zap.Logger, env string) *AuthMiddleware {
	am := &AuthMiddleware{jwtSecret: secret, logger: logger, env: env}
	am.secure = env == "production"

	if env == "production" {
		fmt.Println("this is secure", am.secure )
	}else{
		fmt.Println("this is secure", am.secure )
	}
	return am
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("this is jwt token here here in middleware", am.jwtSecret)
		refreshCookie, err := r.Cookie(RefreshToken)
		if err != nil {
			am.logger.Error("Refresh token cookie missing", zap.Error(err), zap.String("method", r.Method), zap.String("url", r.URL.String()))
			http.Error(w, "unauthorized refresh token", http.StatusUnauthorized)
			return
		}

		claims, err := am.verifyToken(refreshCookie.Value)
		fmt.Println("this is refresh cookie ",refreshCookie.Value)
		if err != nil {
			am.logger.Error("Invalid refresh token",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.String("authHeader", r.Header.Get("Authorization")),
			zap.String("cookie", r.Header.Get("Cookie")), // optional, just for debug
			)
			http.Error(w, "unauthorized access token", http.StatusUnauthorized)
			return
		}

		regenerateAccessToken := false
		accessToken := ""

		if accessCookie, err := r.Cookie(AccessToken); err != nil {
			regenerateAccessToken = true
		} else {
			claims, err = am.verifyToken(accessCookie.Value)
			if err != nil {
				regenerateAccessToken = true
			}
			accessToken = accessCookie.Value
		}

		if regenerateAccessToken {
			accessToken, err = am.GenerateAccessToken(claims.UserID)
			if err != nil {
				am.logger.Error("Error generating access token", zap.Error(err), zap.String("method", r.Method), zap.String("url", r.URL.String()))
				http.Error(w, "error generating access token", http.StatusInternalServerError)
				return
			}
		}

		am.setCookies(w, accessToken, refreshCookie.Value)
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) verifyToken(tokenString string) (*Claims, error) {
	am.logger.Debug("JWT Secret in middleware", zap.String("secret", am.jwtSecret))
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(am.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return claims, fmt.Errorf("invalid token: %w", err)
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return claims, fmt.Errorf("token expired")
	}

	return claims, nil
}

func (am *AuthMiddleware) GenerateAccessToken(userID string) (string, error) {
	return am.generateToken(userID, DefaultAccessTokenDuration)
}

func (am *AuthMiddleware) GenerateRefreshToken(userID string) (string, error) {
	return am.generateToken(userID, DefaultRefreshTokenDuration)
}

func (am *AuthMiddleware) AttachInitialTokens(w http.ResponseWriter, userID string) error {
	am.logger.Debug("JWT Secret in middleware", zap.String("secret", am.jwtSecret))
	accessToken, err := am.GenerateAccessToken(userID)
	if err != nil {
		am.logger.Error("Error generating initial access token", zap.Error(err))
		return err
	}
	refreshToken, err := am.GenerateRefreshToken(userID)
	if err != nil {
		am.logger.Error("Error generating initial refresh token", zap.Error(err))
		return err
	}
	am.setCookies(w, accessToken, refreshToken)
	return nil
}

func (am *AuthMiddleware) generateToken(userID string, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(am.jwtSecret))
}

func (am *AuthMiddleware) setCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessToken,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(DefaultAccessTokenDuration.Seconds()),
		HttpOnly: true,
		Secure:   am.secure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     RefreshToken,
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(DefaultRefreshTokenDuration.Seconds()),
		HttpOnly: true,
		Secure:   am.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetProviderIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
