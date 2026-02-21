package thttp

import (
	"context"
	"net/http"

	"github.com/richardbowden/degrees/internal/services"
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

// SetUserIDInContext is a helper to set user ID in request context
func SetUserIDInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}
