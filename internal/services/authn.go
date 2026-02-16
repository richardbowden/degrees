//go:generate go tool valforge -file $GOFILE
package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardbowden/passwordHash"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/problems"
)

type AuthN struct {
	db dbpg.Storer
}

func NewAuthN(db dbpg.Storer) *AuthN {
	return &AuthN{
		db: db,
	}
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

type Session struct {
	SessionToken string
	UserID       int64
	ExpiresAt    time.Time
}

// VerifyPassword checks if a plaintext password matches the stored hash
func (a *AuthN) VerifyPassword(hashedPassword, plainPassword string) (bool, error) {
	return passwordHash.Validate(plainPassword, hashedPassword)
}

// generateSessionToken creates a cryptographically secure random session token
func (a *AuthN) generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateSession creates a new session for a user
func (a *AuthN) CreateSession(ctx context.Context, userID int64, rememberMe bool, userAgent, ipAddress string) (*Session, error) {
	token, err := a.generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Set expiration based on remember me flag
	var expiresAt time.Time
	if rememberMe {
		expiresAt = time.Now().Add(30 * 24 * time.Hour) // 30 days
	} else {
		expiresAt = time.Now().Add(24 * time.Hour) // 24 hours
	}

	session, err := a.db.CreateSession(ctx, dbpg.CreateSessionParams{
		UserID:       userID,
		SessionToken: token,
		ExpiresAt:    pgtype.Timestamptz{Time: expiresAt, Valid: true},
		UserAgent:    dbpg.StringToPGString(userAgent),
		IpAddress:    dbpg.StringToPGString(ipAddress),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Session{
		SessionToken: session.SessionToken,
		UserID:       session.UserID,
		ExpiresAt:    session.ExpiresAt.Time,
	}, nil
}

// ValidateSession checks if a session token is valid and returns the user ID
func (a *AuthN) ValidateSession(ctx context.Context, sessionToken string) (int64, error) {
	session, err := a.db.GetSessionByToken(ctx, dbpg.GetSessionByTokenParams{SessionToken: sessionToken})
	if err != nil {
		return 0, problems.New(problems.Unauthenticated, "invalid or expired session")
	}

	// Update last activity
	if err := a.db.UpdateSessionActivity(ctx, dbpg.UpdateSessionActivityParams{SessionToken: sessionToken}); err != nil {
		// Log but don't fail - session is still valid
		log := httplog.LogEntry(ctx)
		log.Warn().Err(err).Msg("failed to update session activity")
	}

	return session.UserID, nil
}

// InvalidateSession removes a session
func (a *AuthN) InvalidateSession(ctx context.Context, sessionToken string) error {
	return a.db.DeleteSession(ctx, dbpg.DeleteSessionParams{SessionToken: sessionToken})
}

// InvalidateAllUserSessions removes all sessions for a user
func (a *AuthN) InvalidateAllUserSessions(ctx context.Context, userID int64) error {
	return a.db.DeleteUserSessions(ctx, dbpg.DeleteUserSessionsParams{UserID: userID})
}

// GetUserByEmail retrieves a user by their email address
func (a *AuthN) GetUserByEmail(ctx context.Context, email string) (*dbpg.User, error) {
	user, err := a.db.GetUserByEmail(ctx, dbpg.GetUserByEmailParams{LoginEmail: email})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (a *AuthN) GetUserByID(ctx context.Context, userID int64) (*dbpg.User, error) {
	user, err := a.db.GetUserById(ctx, dbpg.GetUserByIdParams{ID: userID})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *AuthN) HashPassword(pwd1, pwd2 string) (string, error) {
	hashedPassword, err := passwordHash.HashWithDefaults(pwd1, pwd2)
	return hashedPassword, err
}

// UpdateUserPasswordHash updates a user's password hash directly (for admin resets)
func (a *AuthN) UpdateUserPasswordHash(ctx context.Context, userID int64, hashedPassword string) error {
	return a.db.UpdateUserPassword(ctx, dbpg.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: hashedPassword,
	})
}

func (a *AuthN) SetOrUpdatePassword(ctx context.Context, userID int, pwd1, pwd2 string) error {
	hashedPassword, err := a.HashPassword(pwd1, pwd2)
	if err != nil {
		return err
	}
	_ = hashedPassword
	return problems.New(problems.Other, "setorupdate not done yet")

}

// VerifyEmailToken validates a verification token and marks the user's email as verified
func (a *AuthN) VerifyEmailToken(ctx context.Context, token string) error {
	// Get verification record
	verification, err := a.db.GetToken(ctx, dbpg.GetTokenParams{Token: token})
	if err != nil {
		return problems.New(problems.InvalidRequest, "invalid or expired verification token")
	}

	// Update user's sign_up_stage to mark as verified
	_, err = a.db.UpdateUserSignUpStage(ctx, dbpg.UpdateUserSignUpStageParams{
		ID:          verification.UserID,
		SignUpStage: "verified",
	})
	if err != nil {
		return problems.New(problems.Database, "failed to update user verification status", err)
	}

	// Delete the verification token (one-time use)
	err = a.db.DeleteToken(ctx, dbpg.DeleteTokenParams{Token: token})
	if err != nil {
		// Log but don't fail - verification already succeeded
		log := httplog.LogEntry(ctx)
		log.Warn().Err(err).Msg("failed to delete verification token after use")
	}

	return nil
}

// ChangePassword changes the password for an authenticated user
func (a *AuthN) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	log := httplog.LogEntry(ctx)

	// Get user to verify old password
	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return problems.New(problems.Database, "failed to get user", err)
	}

	// Verify old password
	valid, err := a.VerifyPassword(user.PasswordHash, oldPassword)
	if err != nil || !valid {
		return problems.New(problems.InvalidRequest, "old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := passwordHash.HashWithDefaults(newPassword, newPassword)
	if err != nil {
		return problems.New(problems.Internal, "failed to hash new password", err)
	}

	// Update password in database
	err = a.db.UpdateUserPassword(ctx, dbpg.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return problems.New(problems.Database, "failed to update password", err)
	}

	// Invalidate all user sessions for security
	err = a.InvalidateAllUserSessions(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("failed to invalidate sessions after password change")
		// Don't fail - password was changed successfully
	}

	log.Info().Int64("user_id", userID).Msg("password changed successfully")
	return nil
}

// generatePasswordResetToken creates a cryptographically secure random reset token
func (a *AuthN) generatePasswordResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

type PasswordResetResult struct {
	Token string
	Email string
}

// InitiatePasswordReset generates a reset token and stores it for the user
func (a *AuthN) InitiatePasswordReset(ctx context.Context, email string) (*PasswordResetResult, error) {
	log := httplog.LogEntry(ctx)

	// Get user by email
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not for security
		log.Warn().Str("email", email).Msg("password reset requested for non-existent email")
		return nil, nil // Return success but don't send email
	}

	// Delete any existing reset tokens for this user
	err = a.db.DeleteUserPasswordResetTokens(ctx, dbpg.DeleteUserPasswordResetTokensParams{UserID: user.ID})
	if err != nil {
		log.Warn().Err(err).Msg("failed to delete existing reset tokens")
	}

	// Generate reset token
	token, err := a.generatePasswordResetToken()
	if err != nil {
		return nil, problems.New(problems.Internal, "failed to generate reset token", err)
	}

	// Store token with 1 hour expiration
	expiresAt := time.Now().Add(1 * time.Hour)
	err = a.db.CreatePasswordResetToken(ctx, dbpg.CreatePasswordResetTokenParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return nil, problems.New(problems.Database, "failed to create reset token", err)
	}

	log.Info().Int64("user_id", user.ID).Msg("password reset token generated")
	return &PasswordResetResult{
		Token: token,
		Email: user.LoginEmail,
	}, nil
}

// CompletePasswordReset validates reset token and updates password
func (a *AuthN) CompletePasswordReset(ctx context.Context, token, newPassword string) error {
	log := httplog.LogEntry(ctx)

	// Get reset token record (query already checks expiration)
	resetToken, err := a.db.GetPasswordResetToken(ctx, dbpg.GetPasswordResetTokenParams{Token: token})
	if err != nil {
		return problems.New(problems.InvalidRequest, "invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := passwordHash.HashWithDefaults(newPassword, newPassword)
	if err != nil {
		return problems.New(problems.Internal, "failed to hash new password", err)
	}

	// Update password
	err = a.db.UpdateUserPassword(ctx, dbpg.UpdateUserPasswordParams{
		ID:           resetToken.UserID,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return problems.New(problems.Database, "failed to update password", err)
	}

	// Delete the reset token (one-time use)
	err = a.db.DeletePasswordResetToken(ctx, dbpg.DeletePasswordResetTokenParams{Token: token})
	if err != nil {
		log.Warn().Err(err).Msg("failed to delete reset token after use")
	}

	// Invalidate all user sessions for security
	err = a.InvalidateAllUserSessions(ctx, resetToken.UserID)
	if err != nil {
		log.Warn().Err(err).Msg("failed to invalidate sessions after password reset")
	}

	log.Info().Int64("user_id", resetToken.UserID).Msg("password reset completed")
	return nil
}
