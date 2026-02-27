package grpc

import (
	"context"

	fastmail "github.com/richardbowden/degrees/internal/email/genericsmtp"
	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SMTPServiceServer implements the gRPC SMTPService
type SMTPServiceServer struct {
	pb.UnimplementedSMTPServiceServer
	smtpClient *fastmail.Client
	authzSvc   *services.AuthzSvc
}

// NewSMTPServiceServer creates a new SMTP service gRPC server
func NewSMTPServiceServer(smtpClient *fastmail.Client, authzSvc *services.AuthzSvc) *SMTPServiceServer {
	return &SMTPServiceServer{
		smtpClient: smtpClient,
		authzSvc:   authzSvc,
	}
}

// ConfigureSMTP updates the SMTP configuration
func (s *SMTPServiceServer) ConfigureSMTP(ctx context.Context, req *pb.ConfigureSMTPRequest) (*pb.ConfigureSMTPResponse, error) {
	if err := RequireSysop(ctx, s.authzSvc); err != nil {
		return nil, err
	}

	if req.SmtpAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "smtp_address is required")
	}
	if req.SmtpPort == 0 {
		return nil, status.Error(codes.InvalidArgument, "smtp_port is required")
	}
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	err := s.smtpClient.SetConfig(ctx, req.SmtpAddress, int(req.SmtpPort), req.Username, req.Password, req.Identity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to configure SMTP: %v", err)
	}

	return &pb.ConfigureSMTPResponse{
		Success: true,
		Message: "SMTP configuration updated successfully",
	}, nil
}

// GetSMTPStatus returns the current SMTP configuration status
func (s *SMTPServiceServer) GetSMTPStatus(ctx context.Context, req *pb.GetSMTPStatusRequest) (*pb.GetSMTPStatusResponse, error) {
	if err := RequireSysop(ctx, s.authzSvc); err != nil {
		return nil, err
	}

	st := s.smtpClient.GetStatus()
	return &pb.GetSMTPStatusResponse{
		Ready:       st.Ready,
		SmtpAddress: st.SMTPAddress,
		SmtpPort:    int32(st.SMTPPort),
		Username:    st.Username,
		Configured:  st.Configured,
	}, nil
}
