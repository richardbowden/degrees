package repos

import (
	"context"

	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/services"
)

type Catalogue struct {
	store dbpg.Storer
}

func NewCatalogueRepo(store dbpg.Storer) *Catalogue {
	return &Catalogue{
		store: store,
	}
}

func (r *Catalogue) ListCategories(ctx context.Context) ([]dbpg.ServiceCategory, error) {
	return r.store.ListCategories(ctx)
}

func (r *Catalogue) GetCategoryBySlug(ctx context.Context, slug string) (dbpg.ServiceCategory, error) {
	cat, err := r.store.GetCategoryBySlug(ctx, dbpg.GetCategoryBySlugParams{Slug: slug})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.ServiceCategory{}, services.ErrNoRecord
		}
		return dbpg.ServiceCategory{}, err
	}
	return cat, nil
}

func (r *Catalogue) CreateCategory(ctx context.Context, params dbpg.CreateCategoryParams) (dbpg.ServiceCategory, error) {
	return r.store.CreateCategory(ctx, params)
}

func (r *Catalogue) UpdateCategory(ctx context.Context, params dbpg.UpdateCategoryParams) (dbpg.ServiceCategory, error) {
	cat, err := r.store.UpdateCategory(ctx, params)
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.ServiceCategory{}, services.ErrNoRecord
		}
		return dbpg.ServiceCategory{}, err
	}
	return cat, nil
}

func (r *Catalogue) ListServices(ctx context.Context) ([]dbpg.Service, error) {
	return r.store.ListServices(ctx)
}

func (r *Catalogue) ListAllServices(ctx context.Context) ([]dbpg.Service, error) {
	return r.store.ListAllServices(ctx)
}

func (r *Catalogue) ListServicesByCategory(ctx context.Context, categoryID int64) ([]dbpg.Service, error) {
	return r.store.ListServicesByCategory(ctx, dbpg.ListServicesByCategoryParams{CategoryID: categoryID})
}

func (r *Catalogue) GetServiceBySlug(ctx context.Context, slug string) (dbpg.GetServiceBySlugRow, error) {
	svc, err := r.store.GetServiceBySlug(ctx, dbpg.GetServiceBySlugParams{Slug: slug})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.GetServiceBySlugRow{}, services.ErrNoRecord
		}
		return dbpg.GetServiceBySlugRow{}, err
	}
	return svc, nil
}

func (r *Catalogue) GetServiceByID(ctx context.Context, id int64) (dbpg.Service, error) {
	svc, err := r.store.GetServiceByID(ctx, dbpg.GetServiceByIDParams{ID: id})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.Service{}, services.ErrNoRecord
		}
		return dbpg.Service{}, err
	}
	return svc, nil
}

func (r *Catalogue) CreateService(ctx context.Context, params dbpg.CreateServiceParams) (dbpg.Service, error) {
	return r.store.CreateService(ctx, params)
}

func (r *Catalogue) UpdateService(ctx context.Context, params dbpg.UpdateServiceParams) (dbpg.Service, error) {
	svc, err := r.store.UpdateService(ctx, params)
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.Service{}, services.ErrNoRecord
		}
		return dbpg.Service{}, err
	}
	return svc, nil
}

func (r *Catalogue) DeleteService(ctx context.Context, id int64) (dbpg.Service, error) {
	svc, err := r.store.DeleteService(ctx, dbpg.DeleteServiceParams{ID: id})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.Service{}, services.ErrNoRecord
		}
		return dbpg.Service{}, err
	}
	return svc, nil
}

func (r *Catalogue) ListServiceOptions(ctx context.Context, serviceID int64) ([]dbpg.ServiceOption, error) {
	return r.store.ListServiceOptions(ctx, dbpg.ListServiceOptionsParams{ServiceID: serviceID})
}

func (r *Catalogue) ListAllServiceOptions(ctx context.Context, serviceID int64) ([]dbpg.ServiceOption, error) {
	return r.store.ListAllServiceOptions(ctx, dbpg.ListAllServiceOptionsParams{ServiceID: serviceID})
}

func (r *Catalogue) CreateServiceOption(ctx context.Context, params dbpg.CreateServiceOptionParams) (dbpg.ServiceOption, error) {
	return r.store.CreateServiceOption(ctx, params)
}

func (r *Catalogue) UpdateServiceOption(ctx context.Context, params dbpg.UpdateServiceOptionParams) (dbpg.ServiceOption, error) {
	opt, err := r.store.UpdateServiceOption(ctx, params)
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.ServiceOption{}, services.ErrNoRecord
		}
		return dbpg.ServiceOption{}, err
	}
	return opt, nil
}

func (r *Catalogue) DeleteServiceOption(ctx context.Context, id int64) (dbpg.ServiceOption, error) {
	opt, err := r.store.DeleteServiceOption(ctx, dbpg.DeleteServiceOptionParams{ID: id})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.ServiceOption{}, services.ErrNoRecord
		}
		return dbpg.ServiceOption{}, err
	}
	return opt, nil
}

// Vehicle Categories

func (r *Catalogue) ListVehicleCategories(ctx context.Context) ([]dbpg.VehicleCategory, error) {
	return r.store.ListVehicleCategories(ctx)
}

func (r *Catalogue) GetVehicleCategoryByID(ctx context.Context, id int64) (dbpg.VehicleCategory, error) {
	vc, err := r.store.GetVehicleCategoryByID(ctx, dbpg.GetVehicleCategoryByIDParams{ID: id})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.VehicleCategory{}, services.ErrNoRecord
		}
		return dbpg.VehicleCategory{}, err
	}
	return vc, nil
}

func (r *Catalogue) CreateVehicleCategory(ctx context.Context, params dbpg.CreateVehicleCategoryParams) (dbpg.VehicleCategory, error) {
	return r.store.CreateVehicleCategory(ctx, params)
}

func (r *Catalogue) UpdateVehicleCategory(ctx context.Context, params dbpg.UpdateVehicleCategoryParams) (dbpg.VehicleCategory, error) {
	vc, err := r.store.UpdateVehicleCategory(ctx, params)
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.VehicleCategory{}, services.ErrNoRecord
		}
		return dbpg.VehicleCategory{}, err
	}
	return vc, nil
}

func (r *Catalogue) DeleteVehicleCategory(ctx context.Context, id int64) error {
	return r.store.DeleteVehicleCategory(ctx, dbpg.DeleteVehicleCategoryParams{ID: id})
}

// Price Tiers

func (r *Catalogue) ListPriceTiersByService(ctx context.Context, serviceID int64) ([]dbpg.ListPriceTiersByServiceRow, error) {
	return r.store.ListPriceTiersByService(ctx, dbpg.ListPriceTiersByServiceParams{ServiceID: serviceID})
}

func (r *Catalogue) UpsertPriceTier(ctx context.Context, params dbpg.UpsertPriceTierParams) (dbpg.ServicePriceTier, error) {
	return r.store.UpsertPriceTier(ctx, params)
}

func (r *Catalogue) DeletePriceTiersByService(ctx context.Context, serviceID int64) error {
	return r.store.DeletePriceTiersByService(ctx, dbpg.DeletePriceTiersByServiceParams{ServiceID: serviceID})
}

func (r *Catalogue) GetPriceTier(ctx context.Context, serviceID, vehicleCategoryID int64) (dbpg.GetPriceTierRow, error) {
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
