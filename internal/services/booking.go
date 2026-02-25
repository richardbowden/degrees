package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/problems"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, params dbpg.CreateBookingParams) (dbpg.Booking, error)
	CreateBookingService(ctx context.Context, params dbpg.CreateBookingServiceParams) (dbpg.BookingService, error)
	CreateBookingServiceOption(ctx context.Context, params dbpg.CreateBookingServiceOptionParams) (dbpg.BookingServiceOption, error)
	GetBookingByID(ctx context.Context, id int64) (dbpg.GetBookingByIDRow, error)
	ListBookingsByCustomer(ctx context.Context, customerID int64) ([]dbpg.Booking, error)
	ListBookingsByDateRange(ctx context.Context, params dbpg.ListBookingsByDateRangeParams) ([]dbpg.Booking, error)
	ListBookingsForDate(ctx context.Context, params dbpg.ListBookingsForDateParams) ([]dbpg.Booking, error)
	UpdateBookingStatus(ctx context.Context, params dbpg.UpdateBookingStatusParams) (dbpg.Booking, error)
	UpdateBookingPaymentStatus(ctx context.Context, params dbpg.UpdateBookingPaymentStatusParams) (dbpg.Booking, error)
	ListBookingServices(ctx context.Context, bookingID int64) ([]dbpg.ListBookingServicesRow, error)
	ListBookingServiceOptions(ctx context.Context, bookingServiceID int64) ([]dbpg.ListBookingServiceOptionsRow, error)
	GetCartByUserID(ctx context.Context, userID int64) (dbpg.CartSession, error)
	ListCartItems(ctx context.Context, cartSessionID int64) ([]dbpg.ListCartItemsRow, error)
	ClearCart(ctx context.Context, cartSessionID int64) error
	GetCustomerProfileByUserID(ctx context.Context, userID int64) (dbpg.CustomerProfile, error)
	GetServiceByID(ctx context.Context, serviceID int64) (dbpg.Service, error)
	GetVehicleByID(ctx context.Context, vehicleID int64) (Vehicle, error)
	GetPriceTier(ctx context.Context, serviceID, vehicleCategoryID int64) (dbpg.GetPriceTierRow, error)
}

const DepositPercentage = 30

type BookingService struct {
	repo BookingRepository
}

func NewBookingService(repo BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

type CreateBookingFromCartParams struct {
	UserID        int64
	VehicleID     int64
	ScheduledDate string // YYYY-MM-DD
	ScheduledTime string // HH:MM
	Notes         string
}

func (s *BookingService) CreateBookingFromCart(ctx context.Context, params CreateBookingFromCartParams) (*dbpg.Booking, error) {
	// Get customer profile
	customer, err := s.repo.GetCustomerProfileByUserID(ctx, params.UserID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "customer profile not found, please create one first")
		}
		return nil, problems.New(problems.Database, "failed to get customer profile", err)
	}

	// Parse and validate date
	scheduledDate, err := time.Parse("2006-01-02", params.ScheduledDate)
	if err != nil {
		return nil, problems.New(problems.InvalidRequest, "invalid date format, expected YYYY-MM-DD")
	}

	// Enforce minimum 24-hour advance notice
	now := time.Now()
	if scheduledDate.Before(now.Add(24 * time.Hour)) {
		return nil, problems.New(problems.InvalidRequest, "bookings require at least 24 hours advance notice")
	}

	// Parse scheduled time
	scheduledTime, err := time.Parse("15:04", params.ScheduledTime)
	if err != nil {
		return nil, problems.New(problems.InvalidRequest, "invalid time format, expected HH:MM")
	}

	// Get user's cart
	cart, err := s.repo.GetCartByUserID(ctx, params.UserID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "no active cart found")
		}
		return nil, problems.New(problems.Database, "failed to get cart", err)
	}

	// Get cart items
	cartItems, err := s.repo.ListCartItems(ctx, cart.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list cart items", err)
	}
	if len(cartItems) == 0 {
		return nil, problems.New(problems.InvalidRequest, "cart is empty")
	}

	// Resolve vehicle category for tier-based pricing
	var vehicleCategoryID int64
	if params.VehicleID > 0 {
		vehicle, vErr := s.repo.GetVehicleByID(ctx, params.VehicleID)
		if vErr == nil && vehicle.VehicleCategoryID > 0 {
			vehicleCategoryID = vehicle.VehicleCategoryID
		}
	}

	// Calculate totals and estimated duration, using tier price when available
	var subtotal int64
	var totalDuration int32
	for _, item := range cartItems {
		svc, err := s.repo.GetServiceByID(ctx, item.ServiceID)
		if err != nil {
			return nil, problems.New(problems.Database, fmt.Sprintf("failed to get service %d", item.ServiceID), err)
		}
		price := svc.BasePrice
		if vehicleCategoryID > 0 {
			tier, tErr := s.repo.GetPriceTier(ctx, item.ServiceID, vehicleCategoryID)
			if tErr == nil {
				price = tier.Price
			}
		}
		subtotal += price * int64(item.Quantity)
		totalDuration += svc.DurationMinutes * item.Quantity
	}

	totalAmount := subtotal
	depositAmount := totalAmount * DepositPercentage / 100

	pgDate := pgtype.Date{Time: scheduledDate, Valid: true}
	pgTime := pgtype.Time{
		Microseconds: int64(scheduledTime.Hour())*3600000000 + int64(scheduledTime.Minute())*60000000,
		Valid:        true,
	}

	bookingParams := dbpg.CreateBookingParams{
		CustomerID:            customer.ID,
		ScheduledDate:         pgDate,
		ScheduledTime:         pgTime,
		EstimatedDurationMins: totalDuration,
		Status:                dbpg.BookingStatusPendingPayment,
		PaymentStatus:         dbpg.PaymentStatusPending,
		Subtotal:              subtotal,
		DepositAmount:         depositAmount,
		TotalAmount:           totalAmount,
		Notes:                 dbpg.StringToPGString(params.Notes),
	}

	if params.VehicleID > 0 {
		bookingParams.VehicleID = pgtype.Int8{Int64: params.VehicleID, Valid: true}
	}

	booking, err := s.repo.CreateBooking(ctx, bookingParams)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to create booking", err)
	}

	// Snapshot cart items into booking_services with tier-adjusted pricing
	for _, item := range cartItems {
		svc, _ := s.repo.GetServiceByID(ctx, item.ServiceID)
		price := svc.BasePrice
		if vehicleCategoryID > 0 {
			tier, tErr := s.repo.GetPriceTier(ctx, item.ServiceID, vehicleCategoryID)
			if tErr == nil {
				price = tier.Price
			}
		}
		for q := int32(0); q < item.Quantity; q++ {
			_, err := s.repo.CreateBookingService(ctx, dbpg.CreateBookingServiceParams{
				BookingID:      booking.ID,
				ServiceID:      item.ServiceID,
				PriceAtBooking: price,
			})
			if err != nil {
				return nil, problems.New(problems.Database, "failed to create booking service", err)
			}
		}
	}

	// Clear the cart after checkout
	_ = s.repo.ClearCart(ctx, cart.ID)

	return &booking, nil
}

func (s *BookingService) GetBookingByID(ctx context.Context, bookingID int64) (*dbpg.GetBookingByIDRow, error) {
	row, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "booking not found")
		}
		return nil, problems.New(problems.Database, "failed to get booking", err)
	}
	return &row, nil
}

func (s *BookingService) GetMyBooking(ctx context.Context, userID int64, bookingID int64) (*dbpg.GetBookingByIDRow, error) {
	row, err := s.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	// Verify the booking belongs to this user
	if row.CustomerUserID != userID {
		return nil, problems.New(problems.NotExist, "booking not found")
	}

	return row, nil
}

func (s *BookingService) ListMyBookings(ctx context.Context, userID int64) ([]dbpg.Booking, error) {
	customer, err := s.repo.GetCustomerProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return []dbpg.Booking{}, nil
		}
		return nil, problems.New(problems.Database, "failed to get customer profile", err)
	}

	bookings, err := s.repo.ListBookingsByCustomer(ctx, customer.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list bookings", err)
	}
	return bookings, nil
}

func (s *BookingService) ListAllBookings(ctx context.Context, dateFrom, dateTo string) ([]dbpg.Booking, error) {
	fromDate, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return nil, problems.New(problems.InvalidRequest, "invalid date_from format, expected YYYY-MM-DD")
	}
	toDate, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return nil, problems.New(problems.InvalidRequest, "invalid date_to format, expected YYYY-MM-DD")
	}

	bookings, err := s.repo.ListBookingsByDateRange(ctx, dbpg.ListBookingsByDateRangeParams{
		ScheduledDate:   pgtype.Date{Time: fromDate, Valid: true},
		ScheduledDate_2: pgtype.Date{Time: toDate, Valid: true},
	})
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list bookings", err)
	}
	return bookings, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, userID int64, bookingID int64) (*dbpg.Booking, string, error) {
	row, err := s.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, "", err
	}

	if row.CustomerUserID != userID {
		return nil, "", problems.New(problems.NotExist, "booking not found")
	}

	if row.Status == dbpg.BookingStatusCancelled {
		return nil, "", problems.New(problems.InvalidRequest, "booking is already cancelled")
	}

	if row.Status == dbpg.BookingStatusCompleted {
		return nil, "", problems.New(problems.InvalidRequest, "cannot cancel a completed booking")
	}

	booking, err := s.repo.UpdateBookingStatus(ctx, dbpg.UpdateBookingStatusParams{
		ID:     bookingID,
		Status: dbpg.BookingStatusCancelled,
	})
	if err != nil {
		return nil, "", problems.New(problems.Database, "failed to cancel booking", err)
	}

	msg := "booking cancelled"
	if row.ScheduledDate.Valid {
		scheduledDateTime := row.ScheduledDate.Time
		if time.Until(scheduledDateTime) < 24*time.Hour {
			msg = "booking cancelled, deposit may be forfeited due to late cancellation"
		}
	}

	return &booking, msg, nil
}

func (s *BookingService) UpdateBookingStatus(ctx context.Context, bookingID int64, status string) (*dbpg.Booking, error) {
	booking, err := s.repo.UpdateBookingStatus(ctx, dbpg.UpdateBookingStatusParams{
		ID:     bookingID,
		Status: dbpg.BookingStatus(status),
	})
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "booking not found")
		}
		return nil, problems.New(problems.Database, "failed to update booking status", err)
	}
	return &booking, nil
}

func (s *BookingService) CompleteBooking(ctx context.Context, bookingID int64, notes string) (*dbpg.Booking, error) {
	booking, err := s.repo.UpdateBookingStatus(ctx, dbpg.UpdateBookingStatusParams{
		ID:     bookingID,
		Status: dbpg.BookingStatusCompleted,
	})
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "booking not found")
		}
		return nil, problems.New(problems.Database, "failed to complete booking", err)
	}

	// Also mark as fully paid
	booking, err = s.repo.UpdateBookingPaymentStatus(ctx, dbpg.UpdateBookingPaymentStatusParams{
		ID:            bookingID,
		PaymentStatus: dbpg.PaymentStatusFullyPaid,
	})
	if err != nil {
		return nil, problems.New(problems.Database, "failed to update payment status", err)
	}

	return &booking, nil
}

func (s *BookingService) UpdatePaymentStatus(ctx context.Context, bookingID int64, paymentStatus dbpg.PaymentStatus) (*dbpg.Booking, error) {
	booking, err := s.repo.UpdateBookingPaymentStatus(ctx, dbpg.UpdateBookingPaymentStatusParams{
		ID:            bookingID,
		PaymentStatus: paymentStatus,
	})
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "booking not found")
		}
		return nil, problems.New(problems.Database, "failed to update payment status", err)
	}
	return &booking, nil
}

func (s *BookingService) OnDepositPaid(ctx context.Context, bookingID int64) (*dbpg.Booking, error) {
	// Update payment status to deposit_paid
	booking, err := s.repo.UpdateBookingPaymentStatus(ctx, dbpg.UpdateBookingPaymentStatusParams{
		ID:            bookingID,
		PaymentStatus: dbpg.PaymentStatusDepositPaid,
	})
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "booking not found")
		}
		return nil, problems.New(problems.Database, "failed to update payment status", err)
	}

	// Update booking status to deposit_paid
	booking, err = s.repo.UpdateBookingStatus(ctx, dbpg.UpdateBookingStatusParams{
		ID:     bookingID,
		Status: dbpg.BookingStatusDepositPaid,
	})
	if err != nil {
		return nil, problems.New(problems.Database, "failed to update booking status", err)
	}

	return &booking, nil
}

func (s *BookingService) ListBookingServices(ctx context.Context, bookingID int64) ([]dbpg.ListBookingServicesRow, error) {
	return s.repo.ListBookingServices(ctx, bookingID)
}

func (s *BookingService) ListBookingServiceOptions(ctx context.Context, bookingServiceID int64) ([]dbpg.ListBookingServiceOptionsRow, error) {
	return s.repo.ListBookingServiceOptions(ctx, bookingServiceID)
}
