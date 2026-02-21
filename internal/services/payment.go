package services

import (
	"context"
	"errors"

	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/problems"
)

type PaymentBookingRepository interface {
	GetBookingByID(ctx context.Context, id int64) (dbpg.GetBookingByIDRow, error)
	UpdateBookingPaymentStatus(ctx context.Context, params dbpg.UpdateBookingPaymentStatusParams) (dbpg.Booking, error)
	UpdateBookingStatus(ctx context.Context, params dbpg.UpdateBookingStatusParams) (dbpg.Booking, error)
}

type StripeClient interface {
	CreateCheckoutSession(amountCents int64, currency string, bookingID int64, successURL string, cancelURL string) (clientSecret string, err error)
}

type PaymentService struct {
	repo    PaymentBookingRepository
	stripe  StripeClient
	baseURL string
}

func NewPaymentService(repo PaymentBookingRepository, stripe StripeClient, baseURL string) *PaymentService {
	return &PaymentService{
		repo:    repo,
		stripe:  stripe,
		baseURL: baseURL,
	}
}

func (s *PaymentService) CreateDepositSession(ctx context.Context, userID int64, bookingID int64) (string, int64, error) {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return "", 0, problems.New(problems.NotExist, "booking not found")
		}
		return "", 0, problems.New(problems.Database, "failed to get booking", err)
	}

	// Verify the booking belongs to this user
	if booking.CustomerUserID != userID {
		return "", 0, problems.New(problems.NotExist, "booking not found")
	}

	if booking.Status != dbpg.BookingStatusPendingPayment {
		return "", 0, problems.New(problems.InvalidRequest, "booking is not in pending_payment status")
	}

	depositAmount := booking.DepositAmount

	if s.stripe == nil {
		return "", 0, problems.New(problems.Internal, "payment provider not configured")
	}

	successURL := s.baseURL + "/bookings/" + formatInt64(bookingID) + "/success"
	cancelURL := s.baseURL + "/bookings/" + formatInt64(bookingID) + "/cancel"

	clientSecret, err := s.stripe.CreateCheckoutSession(depositAmount, "aud", bookingID, successURL, cancelURL)
	if err != nil {
		return "", 0, problems.New(problems.Internal, "failed to create payment session", err)
	}

	return clientSecret, depositAmount, nil
}

func (s *PaymentService) HandleDepositPaid(ctx context.Context, bookingID int64) (*dbpg.Booking, error) {
	// Update payment status
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

	// Update booking status
	booking, err = s.repo.UpdateBookingStatus(ctx, dbpg.UpdateBookingStatusParams{
		ID:     bookingID,
		Status: dbpg.BookingStatusDepositPaid,
	})
	if err != nil {
		return nil, problems.New(problems.Database, "failed to update booking status", err)
	}

	return &booking, nil
}

func formatInt64(n int64) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	if n < 0 {
		buf = append(buf, '-')
		n = -n
	}
	start := len(buf)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	// reverse digits
	for i, j := start, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
