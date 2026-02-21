package workers

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

const QueueBooking = "booking"

type BookingConfirmationArgs struct {
	BookingID     int64  `json:"booking_id"`
	CustomerEmail string `json:"customer_email"`
	BookingDate   string `json:"booking_date"`
	BookingTime   string `json:"booking_time"`
	DepositAmount int64  `json:"deposit_amount"`
	TotalAmount   int64  `json:"total_amount"`
}

func (BookingConfirmationArgs) Kind() string { return "booking_confirmation" }

func (BookingConfirmationArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: QueueBooking}
}

type BookingNotifier interface {
	SendBookingConfirmation(ctx context.Context, to string, bookingDate string, bookingTime string, depositAmount int64, totalAmount int64) error
}

type BookingConfirmationWorker struct {
	river.WorkerDefaults[BookingConfirmationArgs]
	notifier BookingNotifier
}

func NewBookingConfirmationWorker(notifier BookingNotifier) *BookingConfirmationWorker {
	return &BookingConfirmationWorker{notifier: notifier}
}

func (w *BookingConfirmationWorker) Work(ctx context.Context, job *river.Job[BookingConfirmationArgs]) error {
	log.Info().
		Int64("booking_id", job.Args.BookingID).
		Str("email", job.Args.CustomerEmail).
		Msg("sending booking confirmation")

	if w.notifier == nil {
		return fmt.Errorf("booking notifier not configured - job will retry")
	}

	err := w.notifier.SendBookingConfirmation(
		ctx,
		job.Args.CustomerEmail,
		job.Args.BookingDate,
		job.Args.BookingTime,
		job.Args.DepositAmount,
		job.Args.TotalAmount,
	)
	if err != nil {
		return fmt.Errorf("failed to send booking confirmation: %w", err)
	}

	log.Info().
		Int64("booking_id", job.Args.BookingID).
		Msg("booking confirmation sent")

	return nil
}
