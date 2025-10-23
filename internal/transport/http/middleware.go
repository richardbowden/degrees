package thttp

// func authMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Extract token from header
// 		//token := r.Header.Get("Authorization")
// 		//if token == "" {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte(`{"error":"missing authorization header"}`))
// 		return
// 		//}

// 		//// Validate token and get user ID
// 		//userID, err := s.validateToken(token)
// 		//if err != nil {
// 		//	w.WriteHeader(http.StatusUnauthorized)
// 		//	w.Write([]byte(`{"error":"invalid token"}`))
// 		//	return
// 		//}

// 		// Add to context
// 		//ctx := context.WithValue(r.Context(), "user_id", userID)
// 		//next.ServeHTTP(w, r.WithContext(r.Context()))
// 	})
// }

//
//package middleware
//
//import (
//"context"
//"net/http"
//"strings"
//
//"github.com/typewriterco/p402/internal/problems"
//)
//
//type AuthMiddleware struct {
//	authService AuthService
//}
//
//func NewAuthMiddleware(authSvc AuthService) *AuthMiddleware {
//	return &AuthMiddleware{authService: authSvc}
//}
//
//func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		token := extractToken(r)
//		if token == "" {
//			p := problems.New(problems.Unauthenticated, "authentication required")
//			http.Error(w, p.Error(), p.GetStatus())
//			return
//		}
//
//		user, err := m.authService.ValidateToken(r.Context(), token)
//		if err != nil {
//			p := problems.New(problems.Unauthenticated, "invalid token")
//			http.Error(w, p.Error(), p.GetStatus())
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), "user", user)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
//
//func extractToken(r *http.Request) string {
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
