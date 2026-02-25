package repos

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/services"
)

type Bookings struct {
	store dbpg.Storer
}

func NewBookingRepo(store dbpg.Storer) *Bookings {
	return &Bookings{store: store}
}

func (r *Bookings) CreateBooking(ctx context.Context, params dbpg.CreateBookingParams) (dbpg.Booking, error) {
	tx, err := r.store.GetTX(ctx)
	if err != nil {
		return dbpg.Booking{}, err
	}
	defer tx.Rollback(ctx)

	booking, err := tx.CreateBooking(ctx, params)
	if err != nil {
		return dbpg.Booking{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return dbpg.Booking{}, err
	}

	return booking, nil
}

func (r *Bookings) CreateBookingService(ctx context.Context, params dbpg.CreateBookingServiceParams) (dbpg.BookingService, error) {
	return r.store.CreateBookingService(ctx, params)
}

func (r *Bookings) CreateBookingServiceOption(ctx context.Context, params dbpg.CreateBookingServiceOptionParams) (dbpg.BookingServiceOption, error) {
	return r.store.CreateBookingServiceOption(ctx, params)
}

func (r *Bookings) GetBookingByID(ctx context.Context, id int64) (dbpg.GetBookingByIDRow, error) {
	row, err := r.store.GetBookingByID(ctx, dbpg.GetBookingByIDParams{ID: id})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.GetBookingByIDRow{}, services.ErrNoRecord
		}
		return dbpg.GetBookingByIDRow{}, err
	}
	return row, nil
}

func (r *Bookings) ListBookingsByCustomer(ctx context.Context, customerID int64) ([]dbpg.Booking, error) {
	return r.store.ListBookingsByCustomer(ctx, dbpg.ListBookingsByCustomerParams{CustomerID: customerID})
}

func (r *Bookings) ListBookingsByDateRange(ctx context.Context, params dbpg.ListBookingsByDateRangeParams) ([]dbpg.Booking, error) {
	return r.store.ListBookingsByDateRange(ctx, params)
}

func (r *Bookings) ListBookingsForDate(ctx context.Context, params dbpg.ListBookingsForDateParams) ([]dbpg.Booking, error) {
	return r.store.ListBookingsForDate(ctx, params)
}

func (r *Bookings) UpdateBookingStatus(ctx context.Context, params dbpg.UpdateBookingStatusParams) (dbpg.Booking, error) {
	b, err := r.store.UpdateBookingStatus(ctx, params)
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.Booking{}, services.ErrNoRecord
		}
		return dbpg.Booking{}, err
	}
	return b, nil
}

func (r *Bookings) UpdateBookingPaymentStatus(ctx context.Context, params dbpg.UpdateBookingPaymentStatusParams) (dbpg.Booking, error) {
	b, err := r.store.UpdateBookingPaymentStatus(ctx, params)
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.Booking{}, services.ErrNoRecord
		}
		return dbpg.Booking{}, err
	}
	return b, nil
}

func (r *Bookings) ListBookingServices(ctx context.Context, bookingID int64) ([]dbpg.ListBookingServicesRow, error) {
	return r.store.ListBookingServices(ctx, dbpg.ListBookingServicesParams{BookingID: bookingID})
}

func (r *Bookings) ListBookingServiceOptions(ctx context.Context, bookingServiceID int64) ([]dbpg.ListBookingServiceOptionsRow, error) {
	return r.store.ListBookingServiceOptions(ctx, dbpg.ListBookingServiceOptionsParams{BookingServiceID: bookingServiceID})
}

func (r *Bookings) GetCartByUserID(ctx context.Context, userID int64) (dbpg.CartSession, error) {
	cart, err := r.store.GetCartByUserID(ctx, dbpg.GetCartByUserIDParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.CartSession{}, services.ErrNoRecord
		}
		return dbpg.CartSession{}, err
	}
	return cart, nil
}

func (r *Bookings) ListCartItems(ctx context.Context, cartSessionID int64) ([]dbpg.ListCartItemsRow, error) {
	return r.store.ListCartItems(ctx, dbpg.ListCartItemsParams{CartSessionID: cartSessionID})
}

func (r *Bookings) ClearCart(ctx context.Context, cartSessionID int64) error {
	return r.store.ClearCart(ctx, dbpg.ClearCartParams{CartSessionID: cartSessionID})
}

func (r *Bookings) GetCustomerProfileByUserID(ctx context.Context, userID int64) (dbpg.CustomerProfile, error) {
	cp, err := r.store.GetCustomerProfileByUserID(ctx, dbpg.GetCustomerProfileByUserIDParams{UserID: userID})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.CustomerProfile{}, services.ErrNoRecord
		}
		return dbpg.CustomerProfile{}, err
	}
	return cp, nil
}

func (r *Bookings) GetServiceByID(ctx context.Context, serviceID int64) (dbpg.Service, error) {
	svc, err := r.store.GetServiceByID(ctx, dbpg.GetServiceByIDParams{ID: serviceID})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.Service{}, services.ErrNoRecord
		}
		return dbpg.Service{}, err
	}
	return svc, nil
}

func (r *Bookings) GetVehicleByID(ctx context.Context, vehicleID int64) (services.Vehicle, error) {
	v, err := r.store.GetVehicleByID(ctx, dbpg.GetVehicleByIDParams{ID: vehicleID})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.Vehicle{}, services.ErrNoRecord
		}
		return services.Vehicle{}, err
	}
	return services.Vehicle{
		ID:                v.ID,
		CustomerID:        v.CustomerID,
		VehicleCategoryID: v.VehicleCategoryID.Int64,
	}, nil
}

func (r *Bookings) GetPriceTier(ctx context.Context, serviceID, vehicleCategoryID int64) (dbpg.GetPriceTierRow, error) {
	row, err := r.store.GetPriceTier(ctx, dbpg.GetPriceTierParams{
		ServiceID:         serviceID,
		VehicleCategoryID: vehicleCategoryID,
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.GetPriceTierRow{}, services.ErrNoRecord
		}
		return dbpg.GetPriceTierRow{}, err
	}
	return row, nil
}
