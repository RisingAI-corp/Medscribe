package middleware

import (
	contextLogger "Medscribe/logger"
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/rs/xid"
	"go.uber.org/zap"
)

// Define keys to avoid context key collisions.
type ctxKey string

const (
	ctxKeyCorrelationID ctxKey = "correlation_id"
	ctxKeyIPAddress     ctxKey = "ip_address"
	ctxKeyUserAgent     ctxKey = "user_agent"
	ctxKeyReferrer      ctxKey = "referrer"
	CtxKeyUserID       ctxKey = "user_id"
	// Add more keys as needed.
)

// MetadataMiddleware extracts metadata from the HTTP request and stores it into the request context.
func MetadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Extracting Metadata from the request
		// Generate a unique correlation ID for this request.
		correlationID := xid.New().String()
		ipAddress := getIPAddress(r)
		userAgent := r.Header.Get("User-Agent")
		referrer := r.Header.Get("Referer")
		
		ctx := r.Context()
		var userID string
		if uid, ok := ctx.Value(CtxKeyUserID).(string); ok {
			userID = uid
		} else {
			userID = ""
		}

		// Prepare the context by injecting the metadata.
		
		ctx = context.WithValue(ctx, ctxKeyCorrelationID, correlationID)
		ctx = context.WithValue(ctx, ctxKeyIPAddress, ipAddress)
		ctx = context.WithValue(ctx, ctxKeyUserAgent, userAgent)
		ctx = context.WithValue(ctx, ctxKeyReferrer, referrer)

		logger := contextLogger.Get("prod")
		logger = logger.With(
			zap.String("correlation_id", correlationID),
			zap.String("ip_address", ipAddress),
			zap.String("user_agent", userAgent),
			zap.String("referrer", referrer),
			zap.String("user_id", userID),
		)

		//attaching logger to context and context to request
		ctx = contextLogger.WithCtx(ctx, logger)
		logger .Info("Incoming request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Time("time", time.Now()),
		)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// getIPAddress tries to return the client's real IP address.
func getIPAddress(r *http.Request) string {
	// Check if a reverse proxy or load balancer set the X-Forwarded-For header.
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For can be a comma-separated list of IPs. The first one is typically the real client.
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	// If no proxy is involved, use the remote address.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Fallback: return the original remote address.
	}
	return ip
}
