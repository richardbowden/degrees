package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

type PaymentServiceServer struct {
	pb.UnimplementedPaymentServiceServer
	paymentSvc *services.PaymentService
}

func NewPaymentServer(paymentSvc *services.PaymentService) *PaymentServiceServer {
	return &PaymentServiceServer{
		paymentSvc: paymentSvc,
	}
}

func (s *PaymentServiceServer) CreateDepositSession(ctx context.Context, req *pb.CreateDepositSessionRequest) (*pb.CreateDepositSessionResponse, error) {
	if req.BookingId == 0 {
		return nil, status.Error(codes.InvalidArgument, "booking_id is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	clientSecret, depositAmount, err := s.paymentSvc.CreateDepositSession(ctx, userID, req.BookingId)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.CreateDepositSessionResponse{
		ClientSecret:  clientSecret,
		DepositAmount: depositAmount,
	}, nil
}
