package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

type HistoryServiceServer struct {
	pb.UnimplementedHistoryServiceServer
	historySvc *services.HistoryService
}

func NewHistoryServiceServer(historySvc *services.HistoryService) *HistoryServiceServer {
	return &HistoryServiceServer{
		historySvc: historySvc,
	}
}

func (s *HistoryServiceServer) ListMyHistory(ctx context.Context, req *pb.ListMyHistoryRequest) (*pb.ListMyHistoryResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	details, err := s.historySvc.ListMyHistory(ctx, userID)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	records := make([]*pb.ServiceRecord, len(details))
	for i, d := range details {
		records[i] = serviceRecordDetailToPB(&d)
	}

	return &pb.ListMyHistoryResponse{
		Records: records,
	}, nil
}

func (s *HistoryServiceServer) GetServiceRecord(ctx context.Context, req *pb.GetServiceRecordRequest) (*pb.GetServiceRecordResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "service record id is required")
	}

	detail, err := s.historySvc.GetServiceRecord(ctx, userID, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.GetServiceRecordResponse{
		Record: serviceRecordDetailToPB(detail),
	}, nil
}

func (s *HistoryServiceServer) ListCustomerHistory(ctx context.Context, req *pb.ListCustomerHistoryRequest) (*pb.ListCustomerHistoryResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "customer id is required")
	}

	details, err := s.historySvc.ListCustomerHistory(ctx, userID, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	records := make([]*pb.ServiceRecord, len(details))
	for i, d := range details {
		records[i] = serviceRecordDetailToPB(&d)
	}

	return &pb.ListCustomerHistoryResponse{
		Records: records,
	}, nil
}

func (s *HistoryServiceServer) CreateServiceRecord(ctx context.Context, req *pb.CreateServiceRecordRequest) (*pb.CreateServiceRecordResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.BookingId == 0 {
		return nil, status.Error(codes.InvalidArgument, "booking_id is required")
	}
	if req.CustomerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "customer_id is required")
	}

	var completedDate = req.CompletedDate.AsTime()

	record, err := s.historySvc.CreateServiceRecord(ctx, userID, req.BookingId, req.CustomerId, req.VehicleId, completedDate)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.CreateServiceRecordResponse{
		Record: &pb.ServiceRecord{
			Id:            record.ID,
			BookingId:     record.BookingID,
			CustomerId:    record.CustomerID,
			VehicleId:     record.VehicleID,
			CompletedDate: timestamppb.New(record.CompletedDate),
			CreatedAt:     timestamppb.New(record.CreatedAt),
			UpdatedAt:     timestamppb.New(record.UpdatedAt),
		},
	}, nil
}

func (s *HistoryServiceServer) AddServiceNote(ctx context.Context, req *pb.AddServiceNoteRequest) (*pb.AddServiceNoteResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "service record id is required")
	}
	if req.NoteType == "" {
		return nil, status.Error(codes.InvalidArgument, "note_type is required")
	}
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	note, err := s.historySvc.AddServiceNote(ctx, userID, req.Id, req.NoteType, req.Content, req.IsVisibleToCustomer)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.AddServiceNoteResponse{
		Note: serviceNoteToPB(note),
	}, nil
}

func (s *HistoryServiceServer) AddProductUsed(ctx context.Context, req *pb.AddProductUsedRequest) (*pb.AddProductUsedResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "service record id is required")
	}
	if req.ProductName == "" {
		return nil, status.Error(codes.InvalidArgument, "product_name is required")
	}

	product, err := s.historySvc.AddProductUsed(ctx, userID, req.Id, req.ProductName, req.Notes)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.AddProductUsedResponse{
		Product: &pb.ProductUsed{
			Id:              product.ID,
			ServiceRecordId: product.ServiceRecordID,
			ProductName:     product.ProductName,
			Notes:           product.Notes,
		},
	}, nil
}

func serviceRecordDetailToPB(d *services.ServiceRecordDetail) *pb.ServiceRecord {
	notes := make([]*pb.ServiceNote, len(d.Notes))
	for i, n := range d.Notes {
		notes[i] = serviceNoteToPB(&n)
	}

	products := make([]*pb.ProductUsed, len(d.Products))
	for i, p := range d.Products {
		products[i] = &pb.ProductUsed{
			Id:              p.ID,
			ServiceRecordId: p.ServiceRecordID,
			ProductName:     p.ProductName,
			Notes:           p.Notes,
		}
	}

	photos := make([]*pb.ServicePhoto, len(d.Photos))
	for i, p := range d.Photos {
		photos[i] = &pb.ServicePhoto{
			Id:              p.ID,
			ServiceRecordId: p.ServiceRecordID,
			PhotoType:       p.PhotoType,
			Url:             p.URL,
			Caption:         p.Caption,
			CreatedAt:       timestamppb.New(p.CreatedAt),
		}
	}

	return &pb.ServiceRecord{
		Id:            d.Record.ID,
		BookingId:     d.Record.BookingID,
		CustomerId:    d.Record.CustomerID,
		VehicleId:     d.Record.VehicleID,
		CompletedDate: timestamppb.New(d.Record.CompletedDate),
		Notes:         notes,
		Products:      products,
		Photos:        photos,
		CreatedAt:     timestamppb.New(d.Record.CreatedAt),
		UpdatedAt:     timestamppb.New(d.Record.UpdatedAt),
	}
}

func serviceNoteToPB(n *services.ServiceNote) *pb.ServiceNote {
	return &pb.ServiceNote{
		Id:                  n.ID,
		ServiceRecordId:     n.ServiceRecordID,
		NoteType:            n.NoteType,
		Content:             n.Content,
		IsVisibleToCustomer: n.IsVisibleToCustomer,
		CreatedBy:           n.CreatedBy,
		CreatedAt:           timestamppb.New(n.CreatedAt),
	}
}
