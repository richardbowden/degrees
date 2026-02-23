package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

// UserServiceServer implements the gRPC UserService interface
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userSvc *services.UserService
	authSvc *services.AuthN
}

// NewUserServiceServer creates a new gRPC user service server
func NewUserServiceServer(userSvc *services.UserService, authSvc *services.AuthN) *UserServiceServer {
	return &UserServiceServer{userSvc: userSvc, authSvc: authSvc}
}

// GetUser retrieves a user by ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Users can get their own profile; sysops can get anyone's
	callerID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	if callerID != req.UserId {
		if err := RequireSysop(ctx, s.authSvc); err != nil {
			return nil, err
		}
	}

	user, err := s.userSvc.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:         user.ID,
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			Surname:    user.Surname,
			Email:      user.EMail,
			Enabled:    user.Enabled,
			Sysop:      user.Sysop,
			CreatedOn:  timestamppb.New(user.CreatedOn),
			UpdatedAt:  timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

// UpdateUser updates user profile information
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Users can update their own profile; sysops can update anyone's
	callerID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	if callerID != req.UserId {
		if err := RequireSysop(ctx, s.authSvc); err != nil {
			return nil, err
		}
	}

	updateReq := services.UpdateUserRequest{
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		Surname:    req.Surname,
	}

	user, err := s.userSvc.UpdateUser(ctx, req.UserId, updateReq)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:         user.ID,
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			Surname:    user.Surname,
			Email:      user.EMail,
			Enabled:    user.Enabled,
			Sysop:      user.Sysop,
			CreatedOn:  timestamppb.New(user.CreatedOn),
			UpdatedAt:  timestamppb.New(user.UpdatedAt),
		},
		Message: "user updated successfully",
	}, nil
}

// EnableUser enables a user account
func (s *UserServiceServer) EnableUser(ctx context.Context, req *pb.EnableUserRequest) (*pb.EnableUserResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	err := s.userSvc.EnableUser(ctx, req.UserId)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.EnableUserResponse{
		Message: "user enabled successfully",
	}, nil
}

// DisableUser disables a user account
func (s *UserServiceServer) DisableUser(ctx context.Context, req *pb.DisableUserRequest) (*pb.DisableUserResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	err := s.userSvc.DisableUser(ctx, req.UserId)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.DisableUserResponse{
		Message: "user disabled successfully",
	}, nil
}

// SetUserSysop sets sysop status for a user
func (s *UserServiceServer) SetUserSysop(ctx context.Context, req *pb.SetUserSysopRequest) (*pb.SetUserSysopResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	err := s.userSvc.SetUserSysop(ctx, req.UserId, req.Sysop)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.SetUserSysopResponse{
		Message: "sysop status updated successfully",
	}, nil
}

// ListUsers lists all users
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	users, err := s.userSvc.ListAllUsers(ctx)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	// Convert domain users to protobuf users
	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.User{
			Id:         user.ID,
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			Surname:    user.Surname,
			Email:      user.EMail,
			Enabled:    user.Enabled,
			Sysop:      user.Sysop,
			CreatedOn:  timestamppb.New(user.CreatedOn),
			UpdatedAt:  timestamppb.New(user.UpdatedAt),
		}
	}

	return &pb.ListUsersResponse{Users: pbUsers}, nil
}
