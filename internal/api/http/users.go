package apihttp

// import (
// "context"
// "github.com/danielgtaylor/huma/v2"
// "github.com/rs/zerolog/log"
//
// "net/http"
//
// "github.com/go-chi/httplog"
// "github.com/typewriterco/p402/internal/errs"
// "github.com/typewriterco/p402/internal/services"
// )

// type UserHandler struct {
// 	accSvc *services.UserService
// }
//
// func NewUserHandler(userService *services.UserService) *UserHandler {
// 	ah := &UserHandler{accSvc: userService}
//
// 	return ah
// }

// func (uh *UserHandler) RegisterHandlers(api huma.API) {
// 	log.Info().Msg("Registering user handlers")
// 	huma.Register(api, huma.Operation{
// 		Method:      http.MethodPost,
// 		Path:        "/users/new",
// 		Summary:     "registers a new user",
// 		OperationID: "users-new",
// 	}, uh.NewUser)
//
// 	huma.Register(api, huma.Operation{
// 		Method: http.MethodPost,
// 		Path:   "/users/getaccount/{id}",
//
// 		Summary:     "get account by id",
// 		OperationID: "get-account",
// 	}, uh.GetAccount)
//
// }
//
// type NRequest struct {
// 	Body struct {
// 		services.NewUserRequest
// 	}
// }
//
// type TTT struct {
// 	Body struct {
// 		Name string `json:"name"`
// 	}
// 	ID int `json:"id"`
// }

// func (a *UserHandler) GetAccount(ctx context.Context, params *TTT) (*struct{}, error) {
// 	err := a.accSvc.GetAccount(ctx, params.Body.Name)

//var ee error
//ee = errs.E(errs.Database, "main error", e)
//ee := huma.Error404NotFound("rich error", e)
//l := httplog.LogEntry(ctx)
//l.Error().Err(e).Msg("new logggg")
//eee := fmt.Errorf("failed to get the account")
//ee := problems.O(problems.Validation, "fucked", err)

// return nil, err

//return nil, ee

// }

// func (a *UserHandler) NewUser(ctx context.Context, input *NRequest) (*struct{}, error) {
//
// 	log := httplog.LogEntry(ctx)
//
// 	p, err := services.NewSignUpParams(input.Body.Email, input.Body.FirstName, input.Body.MiddleName, input.Body.Surname, input.Body.Username, input.Body.Password1, input.Body.Password2)
//
// 	if err != nil {
// 		//errs.HTTPErrorResponse(w, log, err)
// 		return nil, huma.Error404NotFound("rich error", err)
// 	}
//
// 	err = a.accSvc.NewUser(ctx, p)
//
// 	if err != nil {
// 		log.Error().Err(err).Msg("")
// 		return nil, err
// 	}
//
// 	log.Info().Str("email", p.Email).Msg("")
//
// 	//w.WriteHeader(http.StatusOK)
// 	return nil, nil
// }

// func (a *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	log := httplog.LogEntry(ctx)
//
// 	err := a.accSvc.Login(ctx)
//
// 	if err != nil {
// 		//handle error here
// 	}
//
// 	errs.HTTPErrorResponse(w, log, errs.E(errs.Invalid, "Login has not been done yet!!!!!!!"))
// }
//
// func (a *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	log := httplog.LogEntry(ctx)
//
// 	err := a.accSvc.Logout(ctx)
//
// 	if err != nil {
// 		//handle error here
// 	}
//
// 	errs.HTTPErrorResponse(w, log, errs.E(errs.Invalid, "Logout has not been done yet!!!!!!!"))
// }
