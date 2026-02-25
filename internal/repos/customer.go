package repos

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/services"
)

type Customer struct {
	store dbpg.Storer
}

func NewCustomerRepo(store dbpg.Storer) *Customer {
	return &Customer{
		store: store,
	}
}

func (r *Customer) CreateProfile(ctx context.Context, userID int64, phone, address, suburb, postcode, notes string) (services.CustomerProfile, error) {
	dbProfile, err := r.store.CreateCustomerProfile(ctx, dbpg.CreateCustomerProfileParams{
		UserID:   userID,
		Phone:    dbpg.StringToPGString(phone),
		Address:  dbpg.StringToPGString(address),
		Suburb:   dbpg.StringToPGString(suburb),
		Postcode: dbpg.StringToPGString(postcode),
		Notes:    dbpg.StringToPGString(notes),
	})
	if err != nil {
		return services.CustomerProfile{}, err
	}
	return dbProfileToService(dbProfile), nil
}

func (r *Customer) GetProfileByUserID(ctx context.Context, userID int64) (services.CustomerProfile, error) {
	dbProfile, err := r.store.GetCustomerProfileByUserID(ctx, dbpg.GetCustomerProfileByUserIDParams{
		UserID: userID,
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.CustomerProfile{}, services.ErrNoRecord
		}
		return services.CustomerProfile{}, err
	}
	return dbProfileToService(dbProfile), nil
}

func (r *Customer) UpdateProfile(ctx context.Context, id int64, phone, address, suburb, postcode, notes string) (services.CustomerProfile, error) {
	dbProfile, err := r.store.UpdateCustomerProfile(ctx, dbpg.UpdateCustomerProfileParams{
		ID:       id,
		Phone:    dbpg.StringToPGString(phone),
		Address:  dbpg.StringToPGString(address),
		Suburb:   dbpg.StringToPGString(suburb),
		Postcode: dbpg.StringToPGString(postcode),
		Notes:    dbpg.StringToPGString(notes),
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.CustomerProfile{}, services.ErrNoRecord
		}
		return services.CustomerProfile{}, err
	}
	return dbProfileToService(dbProfile), nil
}

func (r *Customer) ListCustomers(ctx context.Context, limit, offset int32) ([]services.CustomerProfile, error) {
	dbProfiles, err := r.store.ListCustomers(ctx, dbpg.ListCustomersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	profiles := make([]services.CustomerProfile, len(dbProfiles))
	for i, p := range dbProfiles {
		profiles[i] = dbProfileToService(p)
	}
	return profiles, nil
}

func (r *Customer) CreateVehicle(ctx context.Context, customerID int64, make, model string, year int32, colour, rego, paintType, conditionNotes string, isPrimary bool, vehicleCategoryID int64) (services.Vehicle, error) {
	dbVehicle, err := r.store.CreateVehicle(ctx, dbpg.CreateVehicleParams{
		CustomerID:        customerID,
		Make:              make,
		Model:             model,
		Year:              pgtype.Int4{Int32: year, Valid: year > 0},
		Colour:            dbpg.StringToPGString(colour),
		Rego:              dbpg.StringToPGString(rego),
		PaintType:         dbpg.StringToPGString(paintType),
		ConditionNotes:    dbpg.StringToPGString(conditionNotes),
		IsPrimary:         isPrimary,
		VehicleCategoryID: pgtype.Int8{Int64: vehicleCategoryID, Valid: vehicleCategoryID > 0},
	})
	if err != nil {
		return services.Vehicle{}, err
	}
	return dbVehicleToService(dbVehicle), nil
}

func (r *Customer) GetVehicleByID(ctx context.Context, vehicleID int64) (services.Vehicle, error) {
	dbVehicle, err := r.store.GetVehicleByID(ctx, dbpg.GetVehicleByIDParams{
		ID: vehicleID,
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.Vehicle{}, services.ErrNoRecord
		}
		return services.Vehicle{}, err
	}
	return dbVehicleToService(dbVehicle), nil
}

func (r *Customer) ListVehiclesByCustomer(ctx context.Context, customerID int64) ([]services.Vehicle, error) {
	dbVehicles, err := r.store.ListVehiclesByCustomer(ctx, dbpg.ListVehiclesByCustomerParams{
		CustomerID: customerID,
	})
	if err != nil {
		return nil, err
	}
	vehicles := make([]services.Vehicle, len(dbVehicles))
	for i, v := range dbVehicles {
		vehicles[i] = dbVehicleToService(v)
	}
	return vehicles, nil
}

func (r *Customer) UpdateVehicle(ctx context.Context, id int64, make, model string, year int32, colour, rego, paintType, conditionNotes string, isPrimary bool, vehicleCategoryID int64) (services.Vehicle, error) {
	dbVehicle, err := r.store.UpdateVehicle(ctx, dbpg.UpdateVehicleParams{
		ID:                id,
		Make:              make,
		Model:             model,
		Year:              pgtype.Int4{Int32: year, Valid: year > 0},
		Colour:            dbpg.StringToPGString(colour),
		Rego:              dbpg.StringToPGString(rego),
		PaintType:         dbpg.StringToPGString(paintType),
		ConditionNotes:    dbpg.StringToPGString(conditionNotes),
		IsPrimary:         isPrimary,
		VehicleCategoryID: pgtype.Int8{Int64: vehicleCategoryID, Valid: vehicleCategoryID > 0},
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.Vehicle{}, services.ErrNoRecord
		}
		return services.Vehicle{}, err
	}
	return dbVehicleToService(dbVehicle), nil
}

func (r *Customer) DeleteVehicle(ctx context.Context, vehicleID int64) error {
	return r.store.DeleteVehicle(ctx, dbpg.DeleteVehicleParams{
		ID: vehicleID,
	})
}

func dbProfileToService(p dbpg.CustomerProfile) services.CustomerProfile {
	return services.CustomerProfile{
		ID:        p.ID,
		UserID:    p.UserID,
		Phone:     p.Phone.String,
		Address:   p.Address.String,
		Suburb:    p.Suburb.String,
		Postcode:  p.Postcode.String,
		Notes:     p.Notes.String,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}
}

func dbVehicleToService(v dbpg.Vehicle) services.Vehicle {
	return services.Vehicle{
		ID:                v.ID,
		CustomerID:        v.CustomerID,
		Make:              v.Make,
		Model:             v.Model,
		Year:              v.Year.Int32,
		Colour:            v.Colour.String,
		Rego:              v.Rego.String,
		PaintType:         v.PaintType.String,
		ConditionNotes:    v.ConditionNotes.String,
		IsPrimary:         v.IsPrimary,
		VehicleCategoryID: v.VehicleCategoryID.Int64,
		CreatedAt:         v.CreatedAt.Time,
		UpdatedAt:         v.UpdatedAt.Time,
	}
}
