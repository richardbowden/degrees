//go:generate go tool valforge -file $GOFILE
package services

import (
	"context"

	"time"
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
	FirstName      string `json:"first_name" validate:"required,minlen=2"`
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
}

type NewUserRequest struct {
	Email      string `json:"email" validate:"required,email"`
	FirstName  string `json:"first_name" validate:"required,minlen=2"`
	MiddleName string `json:"middle_name,omitempty,omitzero"`
	Surname    string `json:"surname,omitempty,omitzero"`
	Username   string `json:"username"`
	Password1  string `json:"pwd1" validate:"required"`
	Password2  string `json:"pwd2" validate:"eqfieldsecure=Password1"`
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) (*UserService, error) {
	us := &UserService{repo: repo}
	return us, nil
}

func (us *UserService) GetAccount(ctx context.Context, params string) error {
	//err := problems.O(problems.Database, "error from db call in userService")
	// log := httplog.LogEntry(ctx)
	//
	// if errors.Is(err, ErrNoRecord) {
	// 	eeee := problems.New(problems.NotExist, "account not found")
	//
	// 	detai := problems.Detail{
	// 		Message: "acount_number	",
	// 		Value:   "9999999",
	// 	}
	// 	eeee.AddDetail(detai)
	// 	eeee.AddDetail(fmt.Errorf("test native error"))
	// 	log.Error().Err(err).Msg("base error")
	// 	log.Error().Err(eeee).Msg("problems error")
	// 	return eeee
	// }
	//
	//if errors.Is(err, repos.ErrNoRecord) {
	//	return problems.O(problems.NotExist, "failed to find account")
	//}
	//return errs.E(errs.Database, "UserService.GetAccount")
	//e := errs.E(errs.Validation, "error from service", nil)
	//return e
	return nil
}
