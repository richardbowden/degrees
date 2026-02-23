//go:generate go tool valforge -file $GOFILE
package services

import (
	"context"
	"errors"
	"time"

	"github.com/richardbowden/degrees/internal/problems"
)

type UserState string

const (
	UserStateInitial                  = "Initial"
	UserStateEmailPendingVerification = "EmailPendingVerification"
	UserStateEmailVerified            = "EmailVerified"
	UserStateSignUpComplete           = "SignUpComplete"
	UserStateRejected                 = "SignupRejected"
)

type UserEvent string

const (
	UserEventSubmitSignUp          = "SubmitSignup"
	UserEventClickVerificationLink = "ClickVerificationLink"
	UserEventSignupFailed          = "SignUpFailed"
	UserEventCompleteProfile       = "CompleteProfile"
)

type NewUser struct {
	FirstName      string `json:"first_name" validate:"required,minlen=2"`
	MiddleName     string
	Surname        string
	Username       string
	EMail          string
	State          UserState
	HashedPassword string
}

type User struct {
	ID          int64
	FirstName   string
	MiddleName  string
	Surname     string
	EMail       string
	SignUpStage string
	Enabled     bool
	Sysop       bool
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
	Create(ctx context.Context, params NewUser) (User, error)
	DoesUserExist(ctx context.Context, email string, username string) (bool, bool, error)
	GetUserByID(ctx context.Context, userID int64) (User, error)
	UpdateUser(ctx context.Context, userID int64, firstName string, middleName string, surname string) (User, error)
	UpdateSysop(ctx context.Context, userID int64, sysop bool) error
	UpdateEnabled(ctx context.Context, userID int64, enabled bool) (User, error)
	ListAllUsers(ctx context.Context) ([]User, error)
	IsFirstUser(ctx context.Context) (bool, error)
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
	repo  UserRepository
	ac    *AuthzSvc
	authn *AuthN
}

func NewUserService(repo UserRepository, authz *AuthzSvc, authn *AuthN) (*UserService, error) {
	us := &UserService{repo: repo, ac: authz, authn: authn}
	return us, nil
}

// SetUserSysop updates the sysop flag for a user
func (us *UserService) SetUserSysop(ctx context.Context, userID int64, sysop bool) error {
	return us.repo.UpdateSysop(ctx, userID, sysop)
}

// IsFirstUser returns true if there is exactly one user in the system
func (us *UserService) IsFirstUser(ctx context.Context) (bool, error) {
	return us.repo.IsFirstUser(ctx)
}

func (us *UserService) EmailExists(ctx context.Context, email string) (bool, error) {

	eExists, _, err := us.repo.DoesUserExist(ctx, email, "")

	return eExists, err
}

func (us *UserService) CreateUser(ctx context.Context, nu NewUser) (User, error) {

	user, err := us.repo.Create(ctx, nu)
	return user, err

}

// GetUserByID retrieves a user by ID from the repository
func (us *UserService) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	user, err := us.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "user not found")
		}
		return nil, problems.New(problems.Database, "failed to retrieve user", err)
	}
	return &user, nil
}

func (us *UserService) EnableUser(ctx context.Context, userID int64) error {
	_, err := us.repo.UpdateEnabled(ctx, userID, true)
	return err
}

func (us *UserService) DisableUser(ctx context.Context, userID int64) error {
	_, err := us.repo.UpdateEnabled(ctx, userID, false)
	return err
}

func (us *UserService) ListAllUsers(ctx context.Context) ([]User, error) {
	// Requires new DB query - see repository implementation
	return us.repo.ListAllUsers(ctx)
}

// ResetUserPassword resets a user's password (admin action)
func (us *UserService) ResetUserPassword(ctx context.Context, userID int64, newPassword string) error {
	// Hash the new password (pass same value twice since we don't need confirmation for admin reset)
	hashedPassword, err := us.authn.HashPassword(newPassword, newPassword)
	if err != nil {
		return problems.New(problems.Internal, "failed to hash password", err)
	}

	// Update password directly in database
	err = us.authn.UpdateUserPasswordHash(ctx, userID, hashedPassword)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return problems.New(problems.NotExist, "user not found")
		}
		return problems.New(problems.Database, "failed to update password", err)
	}

	// Invalidate all user sessions for security
	err = us.authn.InvalidateAllUserSessions(ctx, userID)
	if err != nil {
		// Log but don't fail - password was updated successfully
		return nil
	}

	return nil
}

// UpdateUserRequest contains fields that can be updated
type UpdateUserRequest struct {
	FirstName  string
	MiddleName string
	Surname    string
}

// UpdateUser updates user profile information
func (us *UserService) UpdateUser(ctx context.Context, userID int64, req UpdateUserRequest) (*User, error) {
	user, err := us.repo.UpdateUser(ctx, userID, req.FirstName, req.MiddleName, req.Surname)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "user not found")
		}
		return nil, problems.New(problems.Database, "failed to update user", err)
	}
	return &user, nil
}
