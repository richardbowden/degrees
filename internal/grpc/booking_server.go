package grpc

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/richardbowden/degrees/internal/dbpg"
	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

type BookingServiceServer struct {
	pb.UnimplementedBookingServiceServer
	bookingSvc  *services.BookingService
	scheduleSvc *services.ScheduleService
}

func NewBookingServer(bookingSvc *services.BookingService, scheduleSvc *services.ScheduleService) *BookingServiceServer {
	return &BookingServiceServer{
		bookingSvc:  bookingSvc,
		scheduleSvc: scheduleSvc,
	}
}

func (s *BookingServiceServer) CreateBookingFromCart(ctx context.Context, req *pb.CreateBookingFromCartRequest) (*pb.CreateBookingFromCartResponse, error) {
	if req.ScheduledDate == "" {
		return nil, status.Error(codes.InvalidArgument, "scheduled_date is required")
	}
	if req.ScheduledTime == "" {
		return nil, status.Error(codes.InvalidArgument, "scheduled_time is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	booking, err := s.bookingSvc.CreateBookingFromCart(ctx, services.CreateBookingFromCartParams{
		UserID:        userID,
		VehicleID:     req.VehicleId,
		ScheduledDate: req.ScheduledDate,
		ScheduledTime: req.ScheduledTime,
		Notes:         req.Notes,
	})
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.CreateBookingFromCartResponse{
		Booking: bookingToProto(booking),
	}, nil
}

func (s *BookingServiceServer) GetAvailableSlots(ctx context.Context, req *pb.GetAvailableSlotsRequest) (*pb.GetAvailableSlotsResponse, error) {
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}

	slots, err := s.scheduleSvc.GetAvailableSlots(ctx, req.Date, req.DurationMinutes)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbSlots := make([]*pb.AvailableSlot, len(slots))
	for i, slot := range slots {
		pbSlots[i] = &pb.AvailableSlot{
			Date:                  slot.Date,
			Time:                  slot.Time,
			AvailableDurationMins: slot.AvailableDurationMins,
		}
	}

	return &pb.GetAvailableSlotsResponse{
		Slots: pbSlots,
	}, nil
}

func (s *BookingServiceServer) ListMyBookings(ctx context.Context, req *pb.ListMyBookingsRequest) (*pb.ListMyBookingsResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	bookings, err := s.bookingSvc.ListMyBookings(ctx, userID)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbBookings := make([]*pb.Booking, len(bookings))
	for i, b := range bookings {
		pbBookings[i] = dbBookingToProto(&b)
	}

	return &pb.ListMyBookingsResponse{
		Bookings: pbBookings,
	}, nil
}

func (s *BookingServiceServer) GetMyBooking(ctx context.Context, req *pb.GetMyBookingRequest) (*pb.GetMyBookingResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	row, err := s.bookingSvc.GetMyBooking(ctx, userID, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbBooking := bookingRowToProto(row)

	// Enrich with services
	svcs, err := s.bookingSvc.ListBookingServices(ctx, req.Id)
	if err == nil {
		pbBooking.Services = bookingServicesToProto(ctx, s.bookingSvc, svcs)
	}

	return &pb.GetMyBookingResponse{
		Booking: pbBooking,
	}, nil
}

func (s *BookingServiceServer) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	booking, msg, err := s.bookingSvc.CancelBooking(ctx, userID, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.CancelBookingResponse{
		Booking: dbBookingToProto(booking),
		Message: msg,
	}, nil
}

func (s *BookingServiceServer) ListAllBookings(ctx context.Context, req *pb.ListAllBookingsRequest) (*pb.ListAllBookingsResponse, error) {
	if req.DateFrom == "" || req.DateTo == "" {
		return nil, status.Error(codes.InvalidArgument, "date_from and date_to are required")
	}

	bookings, err := s.bookingSvc.ListAllBookings(ctx, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbBookings := make([]*pb.Booking, len(bookings))
	for i, b := range bookings {
		pbBookings[i] = dbBookingToProto(&b)
	}

	return &pb.ListAllBookingsResponse{
		Bookings: pbBookings,
	}, nil
}

func (s *BookingServiceServer) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	row, err := s.bookingSvc.GetBookingByID(ctx, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbBooking := bookingRowToProto(row)

	svcs, err := s.bookingSvc.ListBookingServices(ctx, req.Id)
	if err == nil {
		pbBooking.Services = bookingServicesToProto(ctx, s.bookingSvc, svcs)
	}

	return &pb.GetBookingResponse{
		Booking: pbBooking,
	}, nil
}

func (s *BookingServiceServer) UpdateBookingStatus(ctx context.Context, req *pb.UpdateBookingStatusRequest) (*pb.UpdateBookingStatusResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}

	booking, err := s.bookingSvc.UpdateBookingStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.UpdateBookingStatusResponse{
		Booking: dbBookingToProto(booking),
	}, nil
}

func (s *BookingServiceServer) CompleteBooking(ctx context.Context, req *pb.CompleteBookingRequest) (*pb.CompleteBookingResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	booking, err := s.bookingSvc.CompleteBooking(ctx, req.Id, req.Notes)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.CompleteBookingResponse{
		Booking: dbBookingToProto(booking),
	}, nil
}

// Conversion helpers

func bookingToProto(b *dbpg.Booking) *pb.Booking {
	if b == nil {
		return nil
	}
	return &pb.Booking{
		Id:                    b.ID,
		CustomerId:            b.CustomerID,
		VehicleId:             b.VehicleID.Int64,
		ScheduledDate:         formatPGDate(b.ScheduledDate),
		ScheduledTime:         formatPGTime(b.ScheduledTime),
		EstimatedDurationMins: b.EstimatedDurationMins,
		Status:                string(b.Status),
		PaymentStatus:         string(b.PaymentStatus),
		Subtotal:              b.Subtotal,
		DepositAmount:         b.DepositAmount,
		TotalAmount:           b.TotalAmount,
		Notes:                 b.Notes.String,
		CreatedAt:             timestampFromPG(b.CreatedAt),
		UpdatedAt:             timestampFromPG(b.UpdatedAt),
	}
}

func dbBookingToProto(b *dbpg.Booking) *pb.Booking {
	return bookingToProto(b)
}

func bookingRowToProto(row *dbpg.GetBookingByIDRow) *pb.Booking {
	if row == nil {
		return nil
	}
	pbBooking := &pb.Booking{
		Id:                    row.ID,
		CustomerId:            row.CustomerID,
		VehicleId:             row.VehicleID.Int64,
		ScheduledDate:         formatPGDate(row.ScheduledDate),
		ScheduledTime:         formatPGTime(row.ScheduledTime),
		EstimatedDurationMins: row.EstimatedDurationMins,
		Status:                string(row.Status),
		PaymentStatus:         string(row.PaymentStatus),
		Subtotal:              row.Subtotal,
		DepositAmount:         row.DepositAmount,
		TotalAmount:           row.TotalAmount,
		Notes:                 row.Notes.String,
		CreatedAt:             timestampFromPG(row.CreatedAt),
		UpdatedAt:             timestampFromPG(row.UpdatedAt),
		Customer: &pb.BookingCustomerInfo{
			UserId: row.CustomerUserID,
			Phone:  row.CustomerPhone.String,
		},
	}

	if row.VehicleMake.Valid || row.VehicleModel.Valid || row.VehicleRego.Valid {
		pbBooking.Vehicle = &pb.BookingVehicleInfo{
			Make:  row.VehicleMake.String,
			Model: row.VehicleModel.String,
			Rego:  row.VehicleRego.String,
		}
	}

	return pbBooking
}

func bookingServicesToProto(ctx context.Context, bookingSvc *services.BookingService, svcs []dbpg.ListBookingServicesRow) []*pb.BookingServiceItem {
	items := make([]*pb.BookingServiceItem, len(svcs))
	for i, svc := range svcs {
		item := &pb.BookingServiceItem{
			Id:             svc.ID,
			ServiceId:      svc.ServiceID,
			ServiceName:    svc.ServiceName,
			ServiceSlug:    svc.ServiceSlug,
			PriceAtBooking: svc.PriceAtBooking,
		}

		opts, err := bookingSvc.ListBookingServiceOptions(ctx, svc.ID)
		if err == nil && len(opts) > 0 {
			pbOpts := make([]*pb.BookingServiceOptionItem, len(opts))
			for j, opt := range opts {
				pbOpts[j] = &pb.BookingServiceOptionItem{
					Id:              opt.ID,
					ServiceOptionId: opt.ServiceOptionID,
					OptionName:      opt.OptionName,
					PriceAtBooking:  opt.PriceAtBooking,
				}
			}
			item.Options = pbOpts
		}

		items[i] = item
	}
	return items
}

func formatPGDate(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

func formatPGTime(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}
	totalMins := t.Microseconds / 60000000
	hours := totalMins / 60
	mins := totalMins % 60
	return fmt.Sprintf("%02d:%02d", hours, mins)
}

func timestampFromPG(ts pgtype.Timestamptz) *timestamppb.Timestamp {
	if !ts.Valid {
		return nil
	}
	return timestamppb.New(ts.Time)
}
