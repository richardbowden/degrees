package thttp

import (
	"errors"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/problems"
	"github.com/typewriterco/p402/internal/services"
	"github.com/typewriterco/p402/internal/valgen"
)

type UserHandler struct {
	userSvc *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userSvc: userService}
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

func (uh *UserHandler) NewUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := httplog.LogEntry(ctx)

	user, err := DecodeJSONBody[services.NewUserRequest](r)

	if err != nil {
		log.Error().Err(err).Msg("problem parsing the json body")
		p := problems.New(problems.InvalidRequest, "problem parsing the json body", err)
		problems.WriteHTTPError(w, p)
		return
	}

	err = user.Validate()

	var vErr *valgen.ValidationError
	if errors.As(err, &vErr) {
		p := problems.New(problems.InvalidRequest, "validation errors")

		for _, e := range vErr.Errors {
			p.AddDetail(e)
		}
		problems.WriteHTTPError(w, p)
		return
	}

	err = uh.userSvc.NewUser(ctx, user)

	if err != nil {
		log.Error().Err(err).Msg("")
		problems.WriteHTTPErrorWithErr(w, err)
		return
	}
	//
	//log.Info().Str("email", p.Email).Msg("")
	//w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("List All Users - admin role - protected!"))
}

func (a *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("reset password"))
}

func (a *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("verify email"))
}

func (a *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// log := httplog.LogEntry(ctx)
	//
	// err := a.accSvc.Login(ctx)
	//
	// if err != nil {
	// 	//handle error here
	// }
	//
	// errs.HTTPErrorResponse(w, log, errs.E(errs.Invalid, "Login has not been done yet!!!!!!!"))
	w.Write([]byte("login email"))
}

func (a *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("logout email"))

	// ctx := r.Context()
	// log := httplog.LogEntry(ctx)
	//
	// err := a.accSvc.Logout(ctx)
	//
	// if err != nil {
	// 	//handle error here
	// }
	//
	// errs.HTTPErrorResponse(w, log, errs.E(errs.Invalid, "Logout has not been done yet!!!!!!!"))
}
