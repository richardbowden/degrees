package apihttp

import (
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/errs"
	"github.com/typewriterco/p402/internal/services"
)

type AccountHandler struct {
	accSvc *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	ah := &AccountHandler{accSvc: accountService}

	return ah
}

func (a *AccountHandler) NewAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := httplog.LogEntry(ctx)

	type request struct {
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		MiddleName string `json:"middle_name"`
		Surname    string `json:"surname"`
		Password1  string `json:"password1"`
		Password2  string `json:"password2"`
	}

	req := new(request)

	err := DecodeJSON(w, r, req)
	if err != nil {
		errs.HTTPErrorResponse(w, log, errs.E(errs.InvalidRequest, err))
		return
	}

	p, err := services.NewSignUpParams(req.Email, req.FirstName, req.MiddleName, req.Surname, req.Password1, req.Password2)

	if err != nil {
		errs.HTTPErrorResponse(w, log, err)
		return
	}

	err = a.accSvc.NewAccount(ctx, p)

	if err != nil {
		errs.HTTPErrorResponse(w, log, err)
		return
	}

	log.Info().Str("email", p.Email).Msg("")

	w.WriteHeader(http.StatusOK)
}

func (a *AccountHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := httplog.LogEntry(ctx)

	err := a.accSvc.Login(ctx)

	if err != nil {
		//handle error here
	}

	errs.HTTPErrorResponse(w, log, errs.E(errs.Invalid, "Login has not been done yet!!!!!!!"))
}

func (a *AccountHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := httplog.LogEntry(ctx)

	err := a.accSvc.Logout(ctx)

	if err != nil {
		//handle error here
	}

	errs.HTTPErrorResponse(w, log, errs.E(errs.Invalid, "Logout has not been done yet!!!!!!!"))
}
