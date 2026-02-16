package thttp

import (
	"context"
	"net/http"

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

// RequireAuth middleware ensures the user is authenticated via session
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := httplog.LogEntry(ctx)

		// Get session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Warn().Msg("no session cookie found")
			p := problems.New(problems.Unauthenticated, "authentication required")
			problems.WriteHTTPError(w, p)
			return
		}

		// Validate session
		userID, err := m.authn.ValidateSession(ctx, cookie.Value)
		if err != nil {
			log.Warn().Err(err).Msg("invalid session")
			problems.WriteHTTPErrorWithErr(w, err)
			return
		}

		// Add user ID to context
		ctx = SetUserIDInContext(ctx, userID)

		log.Info().Int64("user_id", userID).Msg("user authenticated via session")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireSysop middleware ensures the user has sysop privileges
func (m *AuthMiddleware) RequireSysop(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := httplog.LogEntry(ctx)

		// Get user ID from context (set by authentication middleware)
		userID, ok := ctx.Value(userIDKey).(int64)
		if !ok {
			log.Warn().Msg("no user ID in context - authentication required")
			p := problems.New(problems.Unauthenticated, "authentication required")
			problems.WriteHTTPError(w, p)
			return
		}

		// Check if user is sysop
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

// RequireSystemAdmin middleware ensures the user has system admin privileges (sysop or admin)
func (m *AuthMiddleware) RequireSystemAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := httplog.LogEntry(ctx)

		// Get user ID from context (set by authentication middleware)
		userID, ok := ctx.Value(userIDKey).(int64)
		if !ok {
			log.Warn().Msg("no user ID in context - authentication required")
			p := problems.New(problems.Unauthenticated, "authentication required")
			problems.WriteHTTPError(w, p)
			return
		}

		// Check if user is system admin
		isAdmin, err := m.authz.IsSystemAdmin(ctx, userID)
		if err != nil {
			log.Error().Err(err).Int64("user_id", userID).Msg("failed to check admin status")
			p := problems.New(problems.Internal, "authorization check failed", err)
			problems.WriteHTTPError(w, p)
			return
		}

		if !isAdmin {
			log.Warn().Int64("user_id", userID).Msg("user attempted to access admin endpoint without privileges")
			p := problems.New(problems.Unauthorized, "system admin privileges required")
			problems.WriteHTTPError(w, p)
			return
		}

		log.Info().Int64("user_id", userID).Msg("system admin access granted")
		next.ServeHTTP(w, r)
	})
}

// SetUserIDInContext is a helper to set user ID in request context (for testing or other auth mechanisms)
func SetUserIDInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}
