package services

import (
	"context"
	"crypto/subtle"
	"errors"
	"github.com/typewriterco/p402/internal/problems"

	"fmt"
	"github.com/mozillazg/go-slugify"

	//"github.com/typewriterco/p402/internal/repos"
	"time"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/errs"
)

type SignUpRequest struct {
	FirstName  string
	MiddleName string
	Surname    string
	Username   string
	Password1  string
	Password2  string
	Email      string
}

type NewAccount struct {
	FirstName      string
	MiddleName     string
	Surname        string
	Username       string
	EMail          string
	SignUpStage    int
	HashedPassword string
}

type Account struct {
	ID          int64
	FirstName   string
	MiddleName  string
	Surname     string
	EMail       string
	SignUpStage int
	Enabled     bool
	CreatedOn   time.Time
	UpdatedAt   time.Time
}

type EmailAddress struct {
	Id         int
	Email      string
	Verified   bool
	VerifiedOn time.Time
	UpdatedOn  time.Time
}

type EmailAddresses []EmailAddress

type UserRepository interface {
	Create(ctx context.Context, params NewAccount) (Account, error)
	DoesUserExist(ctx context.Context, email string, username string) (bool, bool, error)
	PPP(ctx context.Context, email string, username string) error
}

type NewUserRequest struct {
	Email      string `json:"email" format:"email"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty,omitzero"`
	Surname    string `json:"surname,omitempty,omitzero"`
	Username   string `json:"username"`
	Password1  string `json:"pwd1"`
	Password2  string `json:"pwd2"`
}

func NewSignUpParams(email, firstName, middleName, surname, username, pass1, pass2 string) (SignUpRequest, error) {
	f := SignUpRequest{
		Email:      email,
		FirstName:  firstName,
		MiddleName: middleName,
		Username:   username,
		Surname:    surname,
		Password1:  pass1,
		Password2:  pass2,
	}

	if f.FirstName == "" {
		return f, errs.E(errs.Validation, "missing first name", fmt.Errorf("first name needs to be ser"))
	}

	if f.MiddleName != "" && f.Surname == "" {
		return f, fmt.Errorf(
			"middlename is set, but surname is not. If you do not have a middle name, please populate the surname field",
		)
	}

	isAlpha := StartsWithAlpha(f.Username)
	if !isAlpha {
		return f, errs.E(errs.Validation, "username invalid", fmt.Errorf("username must start with a letter"))
	}

	f.Username = slugify.Slugify(f.Username)
	usernameLen := len(f.Username)

	if usernameLen == 0 || usernameLen < 6 || usernameLen > 15 {
		return f, errs.E(errs.Validation, "username length", fmt.Errorf("username needs to be between 6 and 15 letters"))
	}

	if f.Password1 == "" || f.Password2 == "" {
		return f, errs.E(errs.Validation, "Ensure both password fields are populated", fmt.Errorf("a password field has not been populated"))
	}

	if res := subtle.ConstantTimeCompare([]byte(f.Password1), []byte(f.Password2)); res == 0 {
		return f, errs.E(errs.Validation, "Passwords do not match", fmt.Errorf("User supplied passwords do not match"))
	}

	return f, nil
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) (*UserService, error) {
	ac := &UserService{repo: repo}
	return ac, nil
}

func (ac *UserService) GetAccount(ctx context.Context, params string) error {
	//err := problems.O(problems.Database, "error from db call in userService")
	log := httplog.LogEntry(ctx)
	err := ac.repo.PPP(ctx, "email", "username")

	if errors.Is(err, ErrNoRecord) {
		eeee := problems.New(problems.NotExist, "account not found")

		detai := problems.Detail{
			Message: "acount_number	",
			Value:   "9999999",
		}
		eeee.AddDetail(detai)
		eeee.AddDetail(fmt.Errorf("test native error"))
		log.Error().Err(err).Msg("base error")
		log.Error().Err(eeee).Msg("problems error")
		return eeee
	}

	//if errors.Is(err, repos.ErrNoRecord) {
	//	return problems.O(problems.NotExist, "failed to find account")
	//}
	//return errs.E(errs.Database, "UserService.GetAccount")
	//e := errs.E(errs.Validation, "error from service", nil)
	//return e
	return nil
}

func (ac *UserService) NewUser(ctx context.Context, params SignUpRequest) error {

	e, u, err := ac.DoesUserExist(ctx, params.Email, params.Username)

	if err != nil {
		return errs.E(errs.Internal, "failed to check if account exists", err)
	}

	if e && u {
		ee := errs.E(errs.Validation, "validation errors")
		ee.(*errs.Error).AddDetail(params.Email, "email is invalid")
		ee.(*errs.Error).AddDetail(params.Username, "username is invalid")

		return ee

	}

	//TODO(rich): hash password
	hashedPassword := params.Password1

	na := NewAccount{
		FirstName:      params.FirstName,
		MiddleName:     params.MiddleName,
		Surname:        params.Surname,
		Username:       params.Username,
		EMail:          params.Email,
		HashedPassword: hashedPassword,
	}

	account, err := ac.repo.Create(ctx, na)

	_ = account

	if err != nil {
		return err
	}

	return nil
}

func (ac *UserService) Login(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called login.... this does not work yet")

	return nil
}

func (ac *UserService) Logout(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called logout.. this does not work yet")

	return nil
}

func (ac *UserService) DoesUserExist(ctx context.Context, email, username string) (bool, bool, error) {
	l := httplog.LogEntry(ctx)
	l.Debug().Str("subsystem", "accounts").Str("func", "DoesAccountAlreadyExist").Str("email", email).Msg("")

	e, u, err := ac.repo.DoesUserExist(ctx, email, username)

	if err != nil {
		return false, false, err
	}
	return e, u, nil
}
