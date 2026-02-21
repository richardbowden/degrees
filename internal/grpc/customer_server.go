package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

type CustomerServiceServer struct {
	pb.UnimplementedCustomerServiceServer
	customerSvc *services.CustomerService
}

func NewCustomerServiceServer(customerSvc *services.CustomerService) *CustomerServiceServer {
	return &CustomerServiceServer{
		customerSvc: customerSvc,
	}
}

func (s *CustomerServiceServer) GetMyProfile(ctx context.Context, req *pb.GetMyProfileRequest) (*pb.GetMyProfileResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	profile, err := s.customerSvc.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.GetMyProfileResponse{
		Profile: customerProfileToPB(profile),
	}, nil
}

func (s *CustomerServiceServer) UpdateMyProfile(ctx context.Context, req *pb.UpdateMyProfileRequest) (*pb.UpdateMyProfileResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	profile, err := s.customerSvc.UpdateProfile(ctx, userID, req.Phone, req.Address, req.Suburb, req.Postcode)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.UpdateMyProfileResponse{
		Profile: customerProfileToPB(profile),
	}, nil
}

func (s *CustomerServiceServer) ListMyVehicles(ctx context.Context, req *pb.ListMyVehiclesRequest) (*pb.ListMyVehiclesResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	vehicles, err := s.customerSvc.ListVehicles(ctx, userID)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbVehicles := make([]*pb.Vehicle, len(vehicles))
	for i, v := range vehicles {
		pbVehicles[i] = vehicleToPB(&v)
	}

	return &pb.ListMyVehiclesResponse{
		Vehicles: pbVehicles,
	}, nil
}

func (s *CustomerServiceServer) AddVehicle(ctx context.Context, req *pb.AddVehicleRequest) (*pb.AddVehicleResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.Make == "" {
		return nil, status.Error(codes.InvalidArgument, "make is required")
	}
	if req.Model == "" {
		return nil, status.Error(codes.InvalidArgument, "model is required")
	}

	vehicle, err := s.customerSvc.AddVehicle(ctx, userID, req.Make, req.Model, req.Year, req.Colour, req.Rego, req.PaintType, req.ConditionNotes, req.IsPrimary)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.AddVehicleResponse{
		Vehicle: vehicleToPB(vehicle),
	}, nil
}

func (s *CustomerServiceServer) UpdateVehicle(ctx context.Context, req *pb.UpdateVehicleRequest) (*pb.UpdateVehicleResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "vehicle id is required")
	}
	if req.Make == "" {
		return nil, status.Error(codes.InvalidArgument, "make is required")
	}
	if req.Model == "" {
		return nil, status.Error(codes.InvalidArgument, "model is required")
	}

	vehicle, err := s.customerSvc.UpdateVehicle(ctx, userID, req.Id, req.Make, req.Model, req.Year, req.Colour, req.Rego, req.PaintType, req.ConditionNotes, req.IsPrimary)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.UpdateVehicleResponse{
		Vehicle: vehicleToPB(vehicle),
	}, nil
}

func (s *CustomerServiceServer) DeleteVehicle(ctx context.Context, req *pb.DeleteVehicleRequest) (*pb.DeleteVehicleResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "vehicle id is required")
	}

	err := s.customerSvc.DeleteVehicle(ctx, userID, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.DeleteVehicleResponse{
		Success: true,
	}, nil
}

func (s *CustomerServiceServer) ListCustomers(ctx context.Context, req *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	customers, err := s.customerSvc.ListCustomers(ctx, userID, req.Limit, req.Offset)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbCustomers := make([]*pb.CustomerProfile, len(customers))
	for i, c := range customers {
		pbCustomers[i] = customerProfileToPB(&c)
	}

	return &pb.ListCustomersResponse{
		Customers: pbCustomers,
	}, nil
}

func (s *CustomerServiceServer) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "customer id is required")
	}

	profile, vehicles, err := s.customerSvc.GetCustomer(ctx, userID, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbVehicles := make([]*pb.Vehicle, len(vehicles))
	for i, v := range vehicles {
		pbVehicles[i] = vehicleToPB(&v)
	}

	return &pb.GetCustomerResponse{
		Profile:  customerProfileToPB(profile),
		Vehicles: pbVehicles,
	}, nil
}

func customerProfileToPB(p *services.CustomerProfile) *pb.CustomerProfile {
	return &pb.CustomerProfile{
		Id:        p.ID,
		UserId:    p.UserID,
		Phone:     p.Phone,
		Address:   p.Address,
		Suburb:    p.Suburb,
		Postcode:  p.Postcode,
		Notes:     p.Notes,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
}

func vehicleToPB(v *services.Vehicle) *pb.Vehicle {
	return &pb.Vehicle{
		Id:             v.ID,
		CustomerId:     v.CustomerID,
		Make:           v.Make,
		Model:          v.Model,
		Year:           v.Year,
		Colour:         v.Colour,
		Rego:           v.Rego,
		PaintType:      v.PaintType,
		ConditionNotes: v.ConditionNotes,
		IsPrimary:      v.IsPrimary,
		CreatedAt:      timestamppb.New(v.CreatedAt),
		UpdatedAt:      timestamppb.New(v.UpdatedAt),
	}
}
