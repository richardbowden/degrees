package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v5/pgtype"
	notification "github.com/typewriterco/p402/internal/notifications"
	"github.com/typewriterco/p402/internal/problems"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/settings"
)

type SignUp struct {
	user            *UserService
	authn           *AuthN
	authz           *AuthzSvc
	settingsService *settings.Service

	Notifier *notification.Notifier
}

func NewSignUp(userSvc *UserService, authnSvc *AuthN, authzSvc *AuthzSvc, settingsService *settings.Service) *SignUp {
	return &SignUp{
		user:            userSvc,
		authn:           authnSvc,
		authz:           authzSvc,
		settingsService: settingsService,
	}
}

// generateVerificationToken creates a cryptographically secure random verification token
func (s *SignUp) generateVerificationToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *SignUp) Register(ctx context.Context, newUser *NewUserRequest) error {
	log := httplog.LogEntry(ctx)
	exists, err := s.user.EmailExists(ctx, newUser.Email)
	if err != nil {
		return problems.New(problems.Database, "failed to lookup email address", err)
	}

	if exists {
		return problems.New(problems.InvalidRequest, "email address already registered")
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

	if err != nil {
		return problems.New(problems.Internal, "failed to create user", err)
	}

	// Check if this is the first user in the system
	isFirstUser, err := s.user.IsFirstUser(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to check if first user")
		// Continue with signup even if check fails
	} else if isFirstUser {
		log.Info().Int64("user_id", user.ID).Msg("first user detected - granting sysop privileges and auto-verifying")

		// Update database sysop flag
		if err := s.user.SetUserSysop(ctx, user.ID, true); err != nil {
			log.Error().Err(err).Msg("failed to set sysop flag in database for first user")
			// Don't fail signup, but log the error
		}

		// Grant sysop role in FGA
		if err := s.authz.SetUserAsSysop(ctx, user.ID); err != nil {
			log.Error().Err(err).Msg("failed to set sysop role in FGA for first user")
			// Don't fail signup, but log the error
		}

		// Auto-verify first user (no SMTP configured yet)
		_, err = s.authn.db.UpdateUserSignUpStage(ctx, dbpg.UpdateUserSignUpStageParams{
			ID:          user.ID,
			SignUpStage: "verified",
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to auto-verify first user")
		} else {
			log.Info().Int64("user_id", user.ID).Msg("first user auto-verified - can login immediately")
		}

		// Skip email verification for first user
		return nil
	}

	// Check if dev mode is enabled and email verification should be skipped
	devMode := settings.NewDevMode(s.settingsService)
	if devMode.SkipEmailVerification(ctx) {
		log.Warn().Int64("user_id", user.ID).Msg("skipping email verification (dev mode enabled)")

		// Auto-verify user immediately
		_, err = s.authn.db.UpdateUserSignUpStage(ctx, dbpg.UpdateUserSignUpStageParams{
			ID:          user.ID,
			SignUpStage: "verified",
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to auto-verify user in dev mode")
			return problems.New(problems.Database, "failed to verify user", err)
		}

		log.Info().Int64("user_id", user.ID).Msg("user auto-verified in dev mode - can login immediately")
		return nil
	}

	// Generate verification token
	verificationToken, err := s.generateVerificationToken()
	if err != nil {
		return problems.New(problems.Internal, "failed to generate verification token", err)
	}

	// Store verification token in database
	err = s.authn.db.CreateVerificationToken(ctx, dbpg.CreateVerificationTokenParams{
		Token:     verificationToken,
		UserID:    user.ID,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to store verification token")
		return problems.New(problems.Internal, "failed to create verification token", err)
	}

	log.Info().Int64("user_id", user.ID).Str("email", user.EMail).Msg("sending verification email")

	err = s.Notifier.SendEmail(
		ctx,
		notification.TPL_SYSTEM_VERIFY_EMAIL_ADDRESS,
		[]string{user.EMail},
		"Verify Email Address - P402",
		notification.VerifyEmailData{
			EmailVerifyURL: fmt.Sprintf("Please verify your email by using this token: %s", verificationToken),
		},
	)
	if err != nil {
		// Don't fail registration if email sending fails (SMTP might not be configured)
		// Log the token so admin can manually verify the user
		log.Warn().Err(err).
			Int64("user_id", user.ID).
			Str("email", user.EMail).
			Msg("failed to send verification email - SMTP may not be configured. Token is stored in the database.")

		// Registration still succeeds, user just needs manual verification
		return nil
	}

	log.Info().Int64("user_id", user.ID).Msg("verification email sent successfully")
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
