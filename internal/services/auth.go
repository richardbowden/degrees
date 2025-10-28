//go:generate go tool valforge -file $GOFILE
package services

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/richardbowden/passwordHash"
	"github.com/typewriterco/p402/internal/problems"
)

type Auth struct {
	userRepo UserRepository
}

func NewAuth(userRepo UserRepository) *Auth {
	return &Auth{
		userRepo: userRepo,
	}
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

func (s *Auth) Login(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called login.... this does not work yet")

	return nil
}

func (us *Auth) Logout(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called logout.. this does not work yet")

	return nil
}

func (s *Auth) DoesUserExist(ctx context.Context, email, username string) (bool, bool, error) {
	l := httplog.LogEntry(ctx)
	l.Debug().Str("subsystem", "accounts").Str("func", "DoesAccountAlreadyExist").Str("email", email).Msg("")

	e, u, err := s.DoesUserExist(ctx, email, "")

	if err != nil {
		return false, false, err
	}
	return e, u, nil
}

func (s *Auth) Register(ctx context.Context, params *NewUserRequest) error {

	e, _, err := s.userRepo.DoesUserExist(ctx, params.Email, params.Username)

	if err != nil {
		return problems.New(problems.Database, "failed to check if account exists", err)
	}

	if e {
		p := problems.New(problems.InvalidRequest, "email already exists")
		return p
	}

	hashedPassword, err := passwordHash.HashWithDefaults(params.Password1, params.Password2)

	if err != nil {
		return problems.New(problems.Internal, "failed to hash password", err)
	}

	na := NewAccount{
		FirstName:      params.FirstName,
		MiddleName:     params.MiddleName,
		Surname:        params.Surname,
		Username:       params.Username,
		EMail:          params.Email,
		HashedPassword: hashedPassword,
	}

	account, err := s.userRepo.Create(ctx, na)

	_ = account

	if err != nil {
		return err
	}

	return nil
}
