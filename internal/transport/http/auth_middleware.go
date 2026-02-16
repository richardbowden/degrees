package thttp

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/problems"
	"github.com/typewriterco/p402/internal/services"
)

type contextKey string

const userIDKey contextKey = "user_id"

// AuthMiddleware provides authentication and authorization middleware
type AuthMiddleware struct {
	authn *services.AuthN
	authz *services.AuthzSvc
}

func NewAuthMiddleware(authn *services.AuthN, authz *services.AuthzSvc) *AuthMiddleware {
	return &AuthMiddleware{
		authn: authn,
		authz: authz,
	}
}

// CookieToAuthHeader copies session_token cookie into the Authorization header
// if no Authorization header is already present. This bridges cookie-based auth
// (for future web frontends using HttpOnly cookies) to the Bearer token format
// that the gRPC auth interceptor expects.
func (m *AuthMiddleware) CookieToAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			cookie, err := r.Cookie("session_token")
			if err == nil && cookie.Value != "" {
				r.Header.Set("Authorization", "Bearer "+cookie.Value)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RequireBearerAuth validates the Bearer token from the Authorization header.
// Use this for non-gRPC HTTP routes that don't go through the gRPC auth interceptor.
func (m *AuthMiddleware) RequireBearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := httplog.LogEntry(ctx)

		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			log.Warn().Msg("missing or invalid authorization header")
			p := problems.New(problems.Unauthenticated, "authentication required")
			problems.WriteHTTPError(w, p)
			return
		}

		userID, err := m.authn.ValidateSession(ctx, token)
		if err != nil {
			log.Warn().Err(err).Msg("invalid session token")
			problems.WriteHTTPErrorWithErr(w, err)
			return
		}

		ctx = SetUserIDInContext(ctx, userID)
		log.Info().Int64("user_id", userID).Msg("user authenticated via bearer token")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireSysop checks sysop privileges. Must be chained after RequireBearerAuth.
func (m *AuthMiddleware) RequireSysop(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := httplog.LogEntry(ctx)

		userID, ok := ctx.Value(userIDKey).(int64)
		if !ok {
			log.Warn().Msg("no user ID in context - authentication required")
			p := problems.New(problems.Unauthenticated, "authentication required")
			problems.WriteHTTPError(w, p)
			return
		}

		isSysop, err := m.authz.IsSysop(ctx, userID)
		if err != nil {
			log.Error().Err(err).Int64("user_id", userID).Msg("failed to check sysop status")
			p := problems.New(problems.Internal, "authorization check failed", err)
			problems.WriteHTTPError(w, p)
			return
		}

		if !isSysop {
			log.Warn().Int64("user_id", userID).Msg("user attempted to access sysop endpoint without privileges")
			p := problems.New(problems.Unauthorized, "sysop privileges required")
			problems.WriteHTTPError(w, p)
			return
		}

		log.Info().Int64("user_id", userID).Msg("sysop access granted")
		next.ServeHTTP(w, r)
	})
}

// extractBearerToken extracts the token from "Bearer <token>" format.
func extractBearerToken(auth string) string {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

// SetUserIDInContext is a helper to set user ID in request context
func SetUserIDInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}
