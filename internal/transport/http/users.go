package thttp

import (
	"net/http"

	"github.com/typewriterco/p402/internal/accesscontrol"
	"github.com/typewriterco/p402/internal/services"
)

type UserHandler struct {
	userSvc *services.UserService
	ac      *accesscontrol.FGA
}

func NewUserHandler(userService *services.UserService, ac *accesscontrol.FGA) *UserHandler {
	return &UserHandler{userSvc: userService, ac: ac}
}

// func (uh *UserHandler) RegisterPublicRoutes(r chi.Router) {
// 	r.Post("/", uh.NewUser)
// 	//r.Get("/check-username", uh.CheckUsernameAvailability)
// }

// func (uh *UserHandler) RegisterProtectedRoutes(r chi.Router) {
// 	r.Get("/profile", uh.Profile)
// 	//r.Put("/profile", uh.UpdateProfile)
// }

// type NRequest struct {
// 	Body struct {
// 		services.NewUserRequest
// 	}
// }

//type TTT struct {
//	Body struct {
//		Name string `json:"name"`
//	}
//	ID int `json:"id"`
//}

//func (a *UserHandler) GetAccount(ctx context.Context, params *TTT) (*struct{}, error) {
//	err := a.accSvc.GetAccount(ctx, params.Body.Name)

//var ee error
//ee = errs.E(errs.Database, "main error", e)
//ee := huma.Error404NotFound("rich error", e)
//l := httplog.LogEntry(ctx)
//l.Error().Err(e).Msg("new logggg")
//eee := fmt.Errorf("failed to get the account")
//ee := problems.O(problems.Validation, "fucked", err)

//return nil, err

//return nil, ee

// }

func (uh *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User Profile"))
}

func (uh *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("List All Users - admin role - protected!"))
}

func (a *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("reset password"))
}
