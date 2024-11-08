package services

import (
	"context"
	"crypto/subtle"
	"fmt"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/errs"
)

type SignUpParams struct {
	FirstName  string
	Middlename string
	Surname    string
	Password1  string
	Password2  string
	Email      string
}

func NewSignUpParams(email, firstName, middlename, surname, pass1, pass2 string) (SignUpParams, error) {
	f := SignUpParams{
		Email:      email,
		FirstName:  firstName,
		Middlename: middlename,
		Surname:    surname,
		Password1:  pass1,
		Password2:  pass2,
	}

	if f.FirstName == "" {
		return f, errs.E(errs.Validation, "missing first name", fmt.Errorf("first name needs to be ser"))
	}

	if f.Middlename != "" && f.Surname == "" {
		return f, fmt.Errorf(
			"middlename is set, but surname is not. If you do not have a middle name, please populate the surname field",
		)
	}

	if f.Password1 == "" || f.Password2 == "" {
		return f, errs.E(errs.Validation, "Ensure both password fields are populated", fmt.Errorf("a password field has not been populated"))
	}

	if res := subtle.ConstantTimeCompare([]byte(f.Password1), []byte(f.Password2)); res == 0 {
		return f, errs.E(errs.Validation, "Passwords do not match", fmt.Errorf("User supplied passwords do not match"))
	}

	return f, nil
}

type AccountService struct {
	ds dbpg.DataStorer
}

func (ac *AccountService) NewAccount(ctx context.Context, params SignUpParams) error {
	exists, err := ac.DoesAccountAlreadyExist(ctx, params.Email)

	if err != nil {
		return errs.E(errs.Internal, "failed to check if account exists", err)
	}

	if exists {
		return errs.E(errs.Exist, "account already exists", nil)
	}
	return nil
}

func (ac *AccountService) Login(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called login.... this does not work yet")

	return nil
}

func (ac *AccountService) Logout(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called logout.. this does not work yet")

	return nil
}

func (ac *AccountService) DoesAccountAlreadyExist(ctx context.Context, email string) (bool, error) {
	l := httplog.LogEntry(ctx)
	l.Debug().Str("subsystem", "accounts").Str("func", "check_email").Str("email", email).Msg("")
	exists, err := ac.ds.CheckAccountActiveEmailExists(ctx, dbpg.CheckAccountActiveEmailExistsParams{Email: email})

	if err != nil || !exists {
		return false, err
	}

	return true, nil
}

func NewAccountService(ds dbpg.DataStorer) (*AccountService, error) {

	ac := &AccountService{ds: ds}

	return ac, nil
}
