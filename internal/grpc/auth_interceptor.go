package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog/log"
	"github.com/richardbowden/degrees/internal/services"
)

// Context keys for user information
type contextKey string

const (
	UserIDKey     contextKey = "user_id"
	SessionKey    contextKey = "session"
	UserContextKey contextKey = "user"
)

// Public endpoints that don't require authentication
var publicEndpoints = map[string]bool{
	"/degrees.v1.AuthService/Register":               true,
	"/degrees.v1.AuthService/Login":                  true,
	"/degrees.v1.AuthService/VerifyEmail":            true,
	"/degrees.v1.AuthService/ResetPassword":          true,
	"/degrees.v1.AuthService/CompletePasswordReset":  true,
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":       true,

	// Catalogue public read endpoints
	"/degrees.v1.CatalogueService/ListCategories": true,
	"/degrees.v1.CatalogueService/ListServices":   true,
	"/degrees.v1.CatalogueService/GetService":     true,

	// Cart endpoints (supports guest sessions via session token)
	"/degrees.v1.CartService/GetCart":        true,
	"/degrees.v1.CartService/AddCartItem":    true,
	"/degrees.v1.CartService/UpdateCartItem": true,
	"/degrees.v1.CartService/RemoveCartItem": true,
	"/degrees.v1.CartService/ClearCart":      true,

	// Booking public endpoint
	"/degrees.v1.BookingService/GetAvailableSlots": true,
}

// AuthInterceptor creates a gRPC unary interceptor for authentication
func AuthInterceptor(authSvc *services.AuthN) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for public endpoints
		if publicEndpoints[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no metadata provided")
		}

		// Get authorization header
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			log.Warn().Str("method", info.FullMethod).Msg("no authorization header")
			return nil, status.Error(codes.Unauthenticated, "authorization header required")
		}

		// Extract token from "Bearer <token>"
		token := extractBearerToken(authHeaders[0])
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format, expected: Bearer <token>")
		}

		// Validate session
		userID, err := authSvc.ValidateSession(ctx, token)
		if err != nil {
			log.Warn().Err(err).Str("method", info.FullMethod).Msg("invalid session token")
			return nil, status.Error(codes.Unauthenticated, "invalid or expired session token")
		}

		// Add user information to context
		ctx = context.WithValue(ctx, UserIDKey, userID)

		log.Debug().Int64("user_id", userID).Str("method", info.FullMethod).Msg("authenticated request")

		return handler(ctx, req)
	}
}

// extractBearerToken extracts the token from "Bearer <token>" format
func extractBearerToken(auth string) string {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// RequireSysop checks if the user has sysop privileges via FGA
// This can be used in individual handler methods for additional authorization
func RequireSysop(ctx context.Context, authz *services.AuthzSvc) error {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "user not authenticated")
	}

	isSysop, err := authz.IsSysop(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Int64("user_id", userID).Msg("failed to check sysop permission")
		return status.Error(codes.Internal, "failed to verify permissions")
	}

	if !isSysop {
		log.Warn().Int64("user_id", userID).Msg("user attempted sysop-only action without privileges")
		return status.Error(codes.PermissionDenied, "sysop privileges required")
	}

	return nil
}
