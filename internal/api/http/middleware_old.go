package apihttp

// import (
// 	"net/http"

// 	"git.moonshot.com/richard/opcounter/internal/errs"
// 	"git.moonshot.com/richard/opcounter/internal/users"
// 	"github.com/alexedwards/scs/v2"
// 	"github.com/go-chi/httplog"
// )

// const (
// 	CTX_user_id = "user_id"
// )

// func IsAuthed(sessionManager *scs.SessionManager, userService users.UserService) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			ctx := r.Context()

// 			user_id := sessionManager.GetString(ctx, "user_id")

// 			if user_id == "" {
// 				w.Write([]byte("not logged in"))
// 				w.WriteHeader(http.StatusUnauthorized)
// 				return
// 			}

// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// func LogInMiddleware(sessionManager *scs.SessionManager) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			ctx := r.Context()
// 			log := httplog.LogEntry(ctx)

// 			log.Info().Msg("asdasdasdas")

// 			var token string
// 			cookie, err := r.Cookie(sessionManager.Cookie.Name)
// 			if err == nil {
// 				token = cookie.Value
// 			}

// 			ctx, err = sessionManager.Load(r.Context(), token)

// 			if err != nil {
// 				e := errs.E(errs.Internal, "failed to load session", err)
// 				errs.HTTPErrorResponse(w, log, e)
// 				return
// 			}

// 			user, ok := sessionManager.Get(ctx, "user_id").(string)

// 			_ = user

// 			if ok {
// 				http.Redirect(w, r.WithContext(ctx), "/", http.StatusTemporaryRedirect)
// 				return
// 			}

// 			// ses := sessionManager.Get(ctx, "sid")

// 			/*
// 				1. check if there is a valid session
// 				2. check if the user is enabled
// 				3. if 1 and 2 are true, redirect to dashboard or referer
// 				4. if expired or not set then proceed to login handler

// 			*/

// 			// if ses != nil {
// 			// 	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
// 			// 	return
// 			// }

// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }
