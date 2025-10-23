package thttp

//
//package services
//
//import (
//"context"
//"crypto/rand"
//"encoding/base64"
//"time"
//
//"github.com/golang-jwt/jwt/v5"
//"github.com/typewriterco/p402/internal/problems"
//)
//
//type AuthService struct {
//	userRepo    UserRepository
//	jwtSecret   []byte
//	tokenExpiry time.Duration
//}
//
//type AuthClaims struct {
//	UserID   int64    `json:"user_id"`
//	Username string   `json:"username"`
//	Email    string   `json:"email"`
//	Roles    []string `json:"roles"`
//	jwt.RegisteredClaims
//}
//
//type LoginRequest struct {
//	Email    string `json:"email"`
//	Password string `json:"password"`
//}
//
//type LoginResponse struct {
//	AccessToken  string    `json:"access_token"`
//	RefreshToken string    `json:"refresh_token"`
//	ExpiresAt    time.Time `json:"expires_at"`
//	User         UserInfo  `json:"user"`
//}
//
//type UserInfo struct {
//	ID       int64    `json:"id"`
//	Username string   `json:"username"`
//	Email    string   `json:"email"`
//	Roles    []string `json:"roles"`
//}
//
//func NewAuthService(userRepo UserRepository, jwtSecret []byte) *AuthService {
//	return &AuthService{
//		userRepo:    userRepo,
//		jwtSecret:   jwtSecret,
//		tokenExpiry: 24 * time.Hour,
//	}
//}
//
//func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
//	// Validate user credentials
//	user, err := s.userRepo.GetByEmail(ctx, req.Email)
//	if err != nil {
//		return nil, problems.New(problems.Unauthenticated, "invalid credentials")
//	}
//
//	if !s.verifyPassword(user.PasswordHash, req.Password) {
//		return nil, problems.New(problems.Unauthenticated, "invalid credentials")
//	}
//
//	if !user.Enabled {
//		return nil, problems.New(problems.Unauthorized, "account disabled")
//	}
//
//	// Generate tokens
//	accessToken, err := s.generateAccessToken(user)
//	if err != nil {
//		return nil, problems.New(problems.Internal, "failed to generate token")
//	}
//
//	refreshToken, err := s.generateRefreshToken(user.ID)
//	if err != nil {
//		return nil, problems.New(problems.Internal, "failed to generate refresh token")
//	}
//
//	return &LoginResponse{
//		AccessToken:  accessToken,
//		RefreshToken: refreshToken,
//		ExpiresAt:    time.Now().Add(s.tokenExpiry),
//		User: UserInfo{
//			ID:       user.ID,
//			Username: user.Username,
//			Email:    user.Email,
//			Roles:    s.getUserRoles(user.ID), // Implement based on your needs
//		},
//	}, nil
//}
//
//func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*AuthClaims, error) {
//	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return s.jwtSecret, nil
//	})
//
//	if err != nil {
//		return nil, problems.New(problems.Unauthenticated, "invalid token")
//	}
//
//	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
//		return claims, nil
//	}
//
//	return nil, problems.New(problems.Unauthenticated, "invalid token")
//}
//
//func (s *AuthService) generateAccessToken(user User) (string, error) {
//	claims := AuthClaims{
//		UserID:   user.ID,
//		Username: user.Username,
//		Email:    user.Email,
//		Roles:    s.getUserRoles(user.ID),
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
//			IssuedAt:  jwt.NewNumericDate(time.Now()),
//			NotBefore: jwt.NewNumericDate(time.Now()),
//			Issuer:    "p402",
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	return token.SignedString(s.jwtSecret)
//}
//
//package middleware
//
//import (
//"context"
//"net/http"
//"strings"
//
//"github.com/typewriterco/p402/internal/problems"
//"github.com/typewriterco/p402/internal/services"
//)
//
//type contextKey string
//
//const (
//	UserContextKey contextKey = "user"
//)
//
//type AuthMiddleware struct {
//	authService *services.AuthService
//}
//
//func NewAuthMiddleware(authSvc *services.AuthService) *AuthMiddleware {
//	return &AuthMiddleware{authService: authSvc}
//}
//
//func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		token := extractBearerToken(r)
//		if token == "" {
//			writeError(w, problems.New(problems.Unauthenticated, "authentication required"))
//			return
//		}
//
//		claims, err := m.authService.ValidateToken(r.Context(), token)
//		if err != nil {
//			writeError(w, err)
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), UserContextKey, claims)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
//
//func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			claims := GetUserFromContext(r.Context())
//			if claims == nil {
//				writeError(w, problems.New(problems.Unauthenticated, "authentication required"))
//				return
//			}
//
//			if !hasRole(claims.Roles, role) {
//				writeError(w, problems.New(problems.Unauthorized, "insufficient permissions"))
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
//}
//
//func (m *AuthMiddleware) RequireOwnershipOrAdmin(resourceUserID int64) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			claims := GetUserFromContext(r.Context())
//			if claims == nil {
//				writeError(w, problems.New(problems.Unauthenticated, "authentication required"))
//				return
//			}
//
//			// Allow if user owns resource or is admin
//			if claims.UserID != resourceUserID && !hasRole(claims.Roles, "admin") {
//				writeError(w, problems.New(problems.Unauthorized, "access denied"))
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
//}
//
//func GetUserFromContext(ctx context.Context) *services.AuthClaims {
//	if claims, ok := ctx.Value(UserContextKey).(*services.AuthClaims); ok {
//		return claims
//	}
//	return nil
//}
//
//func extractBearerToken(r *http.Request) string {
//	auth := r.Header.Get("Authorization")
//	if auth == "" {
//		return ""
//	}
//
//	parts := strings.Split(auth, " ")
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		return ""
//	}
//
//	return parts[1]
//}
//
//func hasRole(roles []string, role string) bool {
//	for _, r := range roles {
//		if r == role {
//			return true
//		}
//	}
//	return false
//}
//
//func writeError(w http.ResponseWriter, err error) {
//	if p, ok := err.(problems.Problem); ok {
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(p.GetStatus())
//		// Write JSON error response
//		return
//	}
//
//	http.Error(w, err.Error(), http.StatusInternalServerError)
//}
//
//
//
//
//
//
//package handlers
//
//import (
//"context"
//"net/http"
//
//"github.com/danielgtaylor/huma/v2"
//"github.com/typewriterco/p402/internal/services"
//)
//
//type AuthHandler struct {
//	authService *services.AuthService
//	userService *services.UserService
//}
//
//func NewAuthHandler(authSvc *services.AuthService, userSvc *services.UserService) *AuthHandler {
//	return &AuthHandler{
//		authService: authSvc,
//		userService: userSvc,
//	}
//}
//
//func (h *AuthHandler) RegisterRoutes(api huma.API) {
//	// Public routes
//	huma.Register(api, huma.Operation{
//		Method:  http.MethodPost,
//		Path:    "/auth/login",
//		Summary: "Login user",
//		Tags:    []string{"auth"},
//	}, h.Login)
//
//	huma.Register(api, huma.Operation{
//		Method:  http.MethodPost,
//		Path:    "/auth/refresh",
//		Summary: "Refresh token",
//		Tags:    []string{"auth"},
//	}, h.RefreshToken)
//
//	// Protected routes
//	huma.Register(api, huma.Operation{
//		Method:  http.MethodPost,
//		Path:    "/auth/logout",
//		Summary: "Logout user",
//		Tags:    []string{"auth"},
//	}, h.Logout)
//}
//
//type LoginRequest struct {
//	Body services.LoginRequest
//}
//
//type LoginResponse struct {
//	Body services.LoginResponse
//}
//
//func (h *AuthHandler) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
//	resp, err := h.authService.Login(ctx, req.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	return &LoginResponse{Body: *resp}, nil
//}
//
//func (h *AuthHandler) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*LoginResponse, error) {
//	resp, err := h.authService.RefreshToken(ctx, req.Body.RefreshToken)
//	if err != nil {
//		return nil, err
//	}
//
//	return &LoginResponse{Body: *resp}, nil
//}
//
//func (h *AuthHandler) Logout(ctx context.Context, req *struct{}) (*struct{}, error) {
//	// Invalidate token (add to blacklist if needed)
//	return &struct{}{}, nil
//}
//
//
//
