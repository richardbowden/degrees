package repos

import (
	"context"

	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/services"
)

type Schedule struct {
	store dbpg.Storer
}

func NewScheduleRepo(store dbpg.Storer) *Schedule {
	return &Schedule{store: store}
}

func (r *Schedule) GetScheduleConfig(ctx context.Context) ([]dbpg.ScheduleConfig, error) {
	return r.store.GetScheduleConfig(ctx)
}

func (r *Schedule) GetScheduleConfigForDay(ctx context.Context, dayOfWeek int32) (dbpg.ScheduleConfig, error) {
	cfg, err := r.store.GetScheduleConfigForDay(ctx, dbpg.GetScheduleConfigForDayParams{DayOfWeek: dayOfWeek})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.ScheduleConfig{}, services.ErrNoRecord
		}
		return dbpg.ScheduleConfig{}, err
	}
	return cfg, nil
}

func (r *Schedule) UpdateScheduleConfig(ctx context.Context, params dbpg.UpdateScheduleConfigParams) (dbpg.ScheduleConfig, error) {
	return r.store.UpdateScheduleConfig(ctx, params)
}

func (r *Schedule) IsDateBlackedOut(ctx context.Context, params dbpg.IsDateBlackedOutParams) (bool, error) {
	return r.store.IsDateBlackedOut(ctx, params)
}

func (r *Schedule) CreateBlackout(ctx context.Context, params dbpg.CreateBlackoutParams) (dbpg.ScheduleBlackout, error) {
	return r.store.CreateBlackout(ctx, params)
}

func (r *Schedule) DeleteBlackout(ctx context.Context, id int64) error {
	return r.store.DeleteBlackout(ctx, dbpg.DeleteBlackoutParams{ID: id})
}

func (r *Schedule) ListBlackoutDates(ctx context.Context) ([]dbpg.ScheduleBlackout, error) {
	return r.store.ListBlackoutDates(ctx)
}

func (r *Schedule) ListBookingsForDate(ctx context.Context, params dbpg.ListBookingsForDateParams) ([]dbpg.Booking, error) {
	return r.store.ListBookingsForDate(ctx, params)
}
