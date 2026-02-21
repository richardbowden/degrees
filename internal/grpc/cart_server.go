package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/richardbowden/degrees/internal/dbpg"
	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

type CartServiceServer struct {
	pb.UnimplementedCartServiceServer
	cartSvc *services.CartService
}

func NewCartServiceServer(cartSvc *services.CartService) *CartServiceServer {
	return &CartServiceServer{
		cartSvc: cartSvc,
	}
}

func (s *CartServiceServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	userID, sessionToken := s.extractCartIdentity(ctx)

	result, err := s.cartSvc.GetOrCreateCart(ctx, userID, sessionToken)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.GetCartResponse{Cart: cartResultToPB(result)}, nil
}

func (s *CartServiceServer) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest) (*pb.AddCartItemResponse, error) {
	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than 0")
	}

	userID, sessionToken := s.extractCartIdentity(ctx)

	result, err := s.cartSvc.AddItem(ctx, userID, sessionToken, req.ServiceId, req.VehicleId, req.Quantity, req.OptionIds)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.AddCartItemResponse{Cart: cartResultToPB(result)}, nil
}

func (s *CartServiceServer) UpdateCartItem(ctx context.Context, req *pb.UpdateCartItemRequest) (*pb.UpdateCartItemResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than 0")
	}

	userID, sessionToken := s.extractCartIdentity(ctx)

	result, err := s.cartSvc.UpdateItemQuantity(ctx, userID, sessionToken, req.Id, req.Quantity)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.UpdateCartItemResponse{Cart: cartResultToPB(result)}, nil
}

func (s *CartServiceServer) RemoveCartItem(ctx context.Context, req *pb.RemoveCartItemRequest) (*pb.RemoveCartItemResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID, sessionToken := s.extractCartIdentity(ctx)

	result, err := s.cartSvc.RemoveItem(ctx, userID, sessionToken, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.RemoveCartItemResponse{Cart: cartResultToPB(result)}, nil
}

func (s *CartServiceServer) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.ClearCartResponse, error) {
	userID, sessionToken := s.extractCartIdentity(ctx)

	err := s.cartSvc.ClearCart(ctx, userID, sessionToken)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.ClearCartResponse{Success: true}, nil
}

// extractCartIdentity gets the user ID from context (if authenticated)
// or the cart session token from metadata headers.
func (s *CartServiceServer) extractCartIdentity(ctx context.Context) (int64, string) {
	// Try authenticated user first
	userID, ok := GetUserIDFromContext(ctx)
	if ok && userID > 0 {
		return userID, ""
	}

	// Fall back to session token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		tokens := md.Get("x-cart-session")
		if len(tokens) > 0 {
			return 0, tokens[0]
		}
	}

	return 0, ""
}

func cartResultToPB(result *services.CartResult) *pb.Cart {
	cart := &pb.Cart{
		Id:           result.Session.ID,
		SessionToken: result.Session.SessionToken,
		Subtotal:     result.Subtotal,
	}
	if result.Session.ExpiresAt.Valid {
		cart.ExpiresAt = timestamppb.New(result.Session.ExpiresAt.Time)
	}

	cart.Items = make([]*pb.CartItem, len(result.Items))
	for i, item := range result.Items {
		cart.Items[i] = dbCartItemToPB(item)
	}

	return cart
}

func dbCartItemToPB(item dbpg.ListCartItemsRow) *pb.CartItem {
	ci := &pb.CartItem{
		Id:           item.ID,
		ServiceId:    item.ServiceID,
		Quantity:     item.Quantity,
		ServiceName:  item.ServiceName,
		ServicePrice: item.ServicePrice,
	}
	if item.VehicleID.Valid {
		ci.VehicleId = item.VehicleID.Int64
	}
	if item.CreatedAt.Valid {
		ci.CreatedAt = timestamppb.New(item.CreatedAt.Time)
	}
	return ci
}
