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

type ScheduleRepository interface {
	GetScheduleConfig(ctx context.Context) ([]dbpg.ScheduleConfig, error)
	GetScheduleConfigForDay(ctx context.Context, dayOfWeek int32) (dbpg.ScheduleConfig, error)
	UpdateScheduleConfig(ctx context.Context, params dbpg.UpdateScheduleConfigParams) (dbpg.ScheduleConfig, error)
	IsDateBlackedOut(ctx context.Context, params dbpg.IsDateBlackedOutParams) (bool, error)
	CreateBlackout(ctx context.Context, params dbpg.CreateBlackoutParams) (dbpg.ScheduleBlackout, error)
	DeleteBlackout(ctx context.Context, id int64) error
	ListBlackoutDates(ctx context.Context) ([]dbpg.ScheduleBlackout, error)
	ListBookingsForDate(ctx context.Context, params dbpg.ListBookingsForDateParams) ([]dbpg.Booking, error)
}

type AvailableSlot struct {
	Date                  string
	Time                  string
	AvailableDurationMins int32
}

type ScheduleService struct {
	repo ScheduleRepository
}

func NewScheduleService(repo ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) GetScheduleConfig(ctx context.Context) ([]dbpg.ScheduleConfig, error) {
	configs, err := s.repo.GetScheduleConfig(ctx)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to get schedule config", err)
	}
	return configs, nil
}

func (s *ScheduleService) UpdateScheduleConfig(ctx context.Context, params dbpg.UpdateScheduleConfigParams) (*dbpg.ScheduleConfig, error) {
	if params.DayOfWeek < 0 || params.DayOfWeek > 6 {
		return nil, problems.New(problems.InvalidRequest, "day_of_week must be between 0 (Sunday) and 6 (Saturday)")
	}

	cfg, err := s.repo.UpdateScheduleConfig(ctx, params)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to update schedule config", err)
	}
	return &cfg, nil
}

func (s *ScheduleService) AddBlackout(ctx context.Context, dateStr string, reason string) (*dbpg.ScheduleBlackout, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, problems.New(problems.InvalidRequest, "invalid date format, expected YYYY-MM-DD")
	}

	blackout, err := s.repo.CreateBlackout(ctx, dbpg.CreateBlackoutParams{
		Date:   pgtype.Date{Time: date, Valid: true},
		Reason: dbpg.StringToPGString(reason),
	})
	if err != nil {
		return nil, problems.New(problems.Database, "failed to create blackout", err)
	}
	return &blackout, nil
}

func (s *ScheduleService) RemoveBlackout(ctx context.Context, id int64) error {
	err := s.repo.DeleteBlackout(ctx, id)
	if err != nil {
		return problems.New(problems.Database, "failed to remove blackout", err)
	}
	return nil
}

func (s *ScheduleService) GetAvailableSlots(ctx context.Context, dateStr string, durationMinutes int32) ([]AvailableSlot, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, problems.New(problems.InvalidRequest, "invalid date format, expected YYYY-MM-DD")
	}

	if durationMinutes <= 0 {
		durationMinutes = 60 // default to 60 minutes
	}

	// Check if date is blacked out
	isBlackedOut, err := s.repo.IsDateBlackedOut(ctx, dbpg.IsDateBlackedOutParams{
		Date: pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		return nil, problems.New(problems.Database, "failed to check blackout", err)
	}
	if isBlackedOut {
		return []AvailableSlot{}, nil
	}

	// Get schedule config for this day of week (Go: Sunday=0, same as our DB)
	dayOfWeek := int32(date.Weekday())
	config, err := s.repo.GetScheduleConfigForDay(ctx, dayOfWeek)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return []AvailableSlot{}, nil
		}
		return nil, problems.New(problems.Database, "failed to get schedule config", err)
	}

	if !config.IsOpen {
		return []AvailableSlot{}, nil
	}

	// Get existing bookings for this date
	bookings, err := s.repo.ListBookingsForDate(ctx, dbpg.ListBookingsForDateParams{
		ScheduledDate: pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list bookings for date", err)
	}

	// Convert open/close times from microseconds to minutes since midnight
	openMins := int32(config.OpenTime.Microseconds / 60000000)
	closeMins := int32(config.CloseTime.Microseconds / 60000000)
	bufferMins := config.BufferMinutes

	// Build occupied intervals (start, end) in minutes since midnight
	type interval struct {
		start int32
		end   int32
	}
	var occupied []interval
	for _, b := range bookings {
		if !b.ScheduledTime.Valid {
			continue
		}
		bStart := int32(b.ScheduledTime.Microseconds / 60000000)
		bEnd := bStart + b.EstimatedDurationMins + bufferMins
		occupied = append(occupied, interval{start: bStart, end: bEnd})
	}

	// Generate 30-minute slot start times
	var slots []AvailableSlot
	for slotStart := openMins; slotStart+durationMinutes <= closeMins; slotStart += 30 {
		slotEnd := slotStart + durationMinutes

		// Check if this slot conflicts with any occupied interval
		conflict := false
		for _, occ := range occupied {
			// Slots overlap if one starts before the other ends and vice versa
			if slotStart < occ.end && slotEnd+bufferMins > occ.start {
				conflict = true
				break
			}
		}

		if !conflict {
			hours := slotStart / 60
			mins := slotStart % 60
			slots = append(slots, AvailableSlot{
				Date:                  dateStr,
				Time:                  fmt.Sprintf("%02d:%02d", hours, mins),
				AvailableDurationMins: closeMins - slotStart,
			})
		}
	}

	return slots, nil
}
