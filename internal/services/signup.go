package services

import (
	"context"

	"github.com/typewriterco/p402/internal/problems"
)

type SignUp struct {
	user  *UserService
	authn *AuthN
	authz *AuthzSvc
}

func NewSignUp(userSvc *UserService, authnSvc *AuthN, authzSvc *AuthzSvc) *SignUp {
	return &SignUp{
		user:  userSvc,
		authn: authnSvc,
		authz: authzSvc,
	}
}

func (s *SignUp) Register(ctx context.Context, newUser *NewUserRequest) error {

	exists, err := s.user.EmailExists(ctx, newUser.Email)
	if err != nil {
		return problems.New(problems.Database, "failed to lookup email address", err)
	}

	if exists {
		return problems.New(problems.Exist, "email address already registered", nil)
	}

	err = newUser.Validate()
	if err != nil {
		return err
	}

	hashedPassword, err := s.authn.HashPassword(newUser.Password1, newUser.Password2)

	if err != nil {
		return problems.New(problems.Internal, "failed to hash password", err)
	}

	nu := NewUser{
		FirstName:      newUser.FirstName,
		MiddleName:     newUser.MiddleName,
		Surname:        newUser.Surname,
		Username:       newUser.Username,
		EMail:          newUser.Email,
		HashedPassword: hashedPassword,
		State:          UserStateInitial,
	}

	user, err := s.user.CreateUser(ctx, nu)
	_ = user

	return nil
}

// func (s *AuthN) Register(ctx context.Context, params *NewUserRequest) error {
//
// 	e, _, err := s.userRepo.DoesUserExist(ctx, params.Email, params.Username)
//
// 	if err != nil {
// 		return problems.New(problems.Database, "failed to check if account exists", err)
// 	}
//
// 	if e {
// 		p := problems.New(problems.InvalidRequest, "email already exists")
// 		return p
// 	}
//
// 	hashedPassword, err := passwordHash.HashWithDefaults(params.Password1, params.Password2)
//
// 	if err != nil {
// 		return problems.New(problems.Internal, "failed to hash password", err)
// 	}
//
//
// 	account, err := s.userRepo.Create(ctx, na)
//
// 	_ = account
//
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
