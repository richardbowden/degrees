package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notification "github.com/richardbowden/degrees/internal/notifications"
	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

// AuthServiceServer implements the gRPC AuthService interface
type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	authNSvc *services.AuthN
	signUpSvc *services.SignUp
	baseURL   string
}

// NewAuthServiceServer creates a new gRPC auth service server
func NewAuthServiceServer(authNSvc *services.AuthN, signUpSvc *services.SignUp, baseURL string) *AuthServiceServer {
	return &AuthServiceServer{
		authNSvc:  authNSvc,
		signUpSvc: signUpSvc,
		baseURL:   baseURL,
	}
}

// Register registers a new user
func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Validate required fields
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.FirstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first_name is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.Password != req.PasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "passwords do not match")
	}

	// Create NewUserRequest
	newUser := &services.NewUserRequest{
		Email:      req.Email,
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		Surname:    req.Surname,
		Username:   req.Username,
		Password1:  req.Password,
		Password2:  req.PasswordConfirm,
	}

	// Call signup service
	err := s.signUpSvc.Register(ctx, newUser)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.RegisterResponse{
		Message: "registration successful, please check your email to verify your account",
	}, nil
}

// VerifyEmail verifies a user's email address
func (s *AuthServiceServer) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Verify the email token
	err := s.authNSvc.VerifyEmailToken(ctx, req.Token)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.VerifyEmailResponse{
		Message: "email verified successfully",
	}, nil
}

// Login authenticates a user
func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Get user by email
	user, err := s.authNSvc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Verify password
	valid, err := s.authNSvc.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil || !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Check if user is enabled
	if !user.Enabled {
		return nil, status.Error(codes.PermissionDenied, "account is disabled")
	}

	// Check if user has verified their email
	if user.SignUpStage != "verified" {
		return nil, status.Error(codes.PermissionDenied, "please verify your email before logging in")
	}

	// Create session (no UserAgent or RemoteAddr in gRPC context for now)
	session, err := s.authNSvc.CreateSession(ctx, user.ID, false, "", "")
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.LoginResponse{
		SessionToken: session.SessionToken,
		User: &pb.User{
			Id:         user.ID,
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName.String,
			Surname:    user.Surname.String,
			Email:      user.LoginEmail,
			Enabled:    user.Enabled,
		},
		Message: "login successful",
	}, nil
}

// Logout logs out a user
func (s *AuthServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.SessionToken == "" {
		return nil, status.Error(codes.InvalidArgument, "session_token is required")
	}

	// Invalidate session
	err := s.authNSvc.InvalidateSession(ctx, req.SessionToken)
	if err != nil {
		// Log but don't fail - session might already be invalid
		return &pb.LogoutResponse{
			Message: "logout completed",
		}, nil
	}

	return &pb.LogoutResponse{
		Message: "logout successful",
	}, nil
}

// ChangePassword changes a user's password
func (s *AuthServiceServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if req.OldPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "old_password is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}
	if req.NewPassword != req.NewPasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "new passwords do not match")
	}

	// Get authenticated user ID from context
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Change password
	err := s.authNSvc.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.ChangePasswordResponse{
		Message: "password changed successfully",
	}, nil
}

// ResetPassword initiates password reset
func (s *AuthServiceServer) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Initiate password reset (generates token and returns it for email)
	result, err := s.authNSvc.InitiatePasswordReset(ctx, req.Email)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	// Send password reset email if token was generated
	if result != nil && s.signUpSvc != nil && s.signUpSvc.Notifier != nil {
		resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.baseURL, result.Token)
		err := s.signUpSvc.Notifier.SendEmail(
			ctx,
			notification.TPL_SYSTEM_PASSWORD_RESET,
			[]string{result.Email},
			"Password Reset Request",
			notification.PasswordResetData{
				ResetLink: resetLink,
				Email:     result.Email,
			},
		)
		if err != nil {
			// Log error but don't fail the request to prevent email enumeration
			// The token is already stored, user can still use it if they have it
		}
	}

	// Always return success to prevent email enumeration
	return &pb.ResetPasswordResponse{
		Message: "if the email exists, a password reset link has been sent",
	}, nil
}

// CompletePasswordReset completes password reset with token
func (s *AuthServiceServer) CompletePasswordReset(ctx context.Context, req *pb.CompletePasswordResetRequest) (*pb.CompletePasswordResetResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}
	if req.NewPassword != req.NewPasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "passwords do not match")
	}

	// Complete password reset
	err := s.authNSvc.CompletePasswordReset(ctx, req.Token, req.NewPassword)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.CompletePasswordResetResponse{
		Message: "password reset successfully",
	}, nil
}
