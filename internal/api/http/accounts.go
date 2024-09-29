package apihttp

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/go-playground/validator/v10"
	"github.com/typewriterco/p402/internal/errs"
	"github.com/typewriterco/p402/internal/services"
)

// this is prob slow...
var val *validator.Validate

func init() {
	val = validator.New(validator.WithRequiredStructEnabled())
}

type AccountHandler struct {
	accSvc *services.AccountService
}

func NewAccountHandler(accountServuce *services.AccountService) *AccountHandler {
	ah := &AccountHandler{accSvc: accountServuce}

	return ah
}

type signUpParams struct {
	FirstName, Middlename, Surname string
	Password1, Password2           string
	Email                          string
}

func NewSignUpParams(email, firstName, middlename, surname, pass1, pass2 string) (signUpParams, error) {
	f := signUpParams{
		Email:      email,
		FirstName:  firstName,
		Middlename: middlename,
		Surname:    surname,
		Password1:  pass1,
		Password2:  pass2,
	}

	if f.FirstName == "" {
		return f, fmt.Errorf("first name needs to be ser")
	}

	if f.Middlename != "" && f.Surname == "" {
		return f, fmt.Errorf(
			"middlename is set, but surname is not. If you do not have a middle name, please populate the surname field",
		)
	}

	if res := subtle.ConstantTimeCompare([]byte(f.Password1), []byte(f.Password2)); res == 0 {
		return f, fmt.Errorf("passwords do not match")
	}

	return f, nil
}

func (a *AccountHandler) SignUp(w http.ResponseWriter, r *http.Request) {
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

	p, err := NewSignUpParams(req.Email, req.FirstName, req.MiddleName, req.Surname, req.Password1, req.Password2)

	if err != nil {
		errs.HTTPErrorResponse(w, log, err)
		return
	}

	err = a.accSvc.NewAccount(ctx)

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
