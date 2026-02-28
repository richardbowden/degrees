package services

import (
	"context"
	"errors"

	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/problems"
)

type CatalogueRepository interface {
	ListCategories(ctx context.Context) ([]dbpg.ServiceCategory, error)
	GetCategoryBySlug(ctx context.Context, slug string) (dbpg.ServiceCategory, error)
	CreateCategory(ctx context.Context, params dbpg.CreateCategoryParams) (dbpg.ServiceCategory, error)
	UpdateCategory(ctx context.Context, params dbpg.UpdateCategoryParams) (dbpg.ServiceCategory, error)
	ListServices(ctx context.Context) ([]dbpg.Service, error)
	ListAllServices(ctx context.Context) ([]dbpg.Service, error)
	ListServicesByCategory(ctx context.Context, categoryID int64) ([]dbpg.Service, error)
	GetServiceBySlug(ctx context.Context, slug string) (dbpg.GetServiceBySlugRow, error)
	GetServiceByID(ctx context.Context, id int64) (dbpg.Service, error)
	CreateService(ctx context.Context, params dbpg.CreateServiceParams) (dbpg.Service, error)
	UpdateService(ctx context.Context, params dbpg.UpdateServiceParams) (dbpg.Service, error)
	DeleteService(ctx context.Context, id int64) (dbpg.Service, error)
	ListServiceOptions(ctx context.Context, serviceID int64) ([]dbpg.ServiceOption, error)
	ListAllServiceOptions(ctx context.Context, serviceID int64) ([]dbpg.ServiceOption, error)
	CreateServiceOption(ctx context.Context, params dbpg.CreateServiceOptionParams) (dbpg.ServiceOption, error)
	UpdateServiceOption(ctx context.Context, params dbpg.UpdateServiceOptionParams) (dbpg.ServiceOption, error)
	DeleteServiceOption(ctx context.Context, id int64) (dbpg.ServiceOption, error)
	ListVehicleCategories(ctx context.Context) ([]dbpg.VehicleCategory, error)
	GetVehicleCategoryByID(ctx context.Context, id int64) (dbpg.VehicleCategory, error)
	CreateVehicleCategory(ctx context.Context, params dbpg.CreateVehicleCategoryParams) (dbpg.VehicleCategory, error)
	UpdateVehicleCategory(ctx context.Context, params dbpg.UpdateVehicleCategoryParams) (dbpg.VehicleCategory, error)
	DeleteVehicleCategory(ctx context.Context, id int64) error
	ListPriceTiersByService(ctx context.Context, serviceID int64) ([]dbpg.ListPriceTiersByServiceRow, error)
	UpsertPriceTier(ctx context.Context, params dbpg.UpsertPriceTierParams) (dbpg.ServicePriceTier, error)
	DeletePriceTiersByService(ctx context.Context, serviceID int64) error
	GetPriceTier(ctx context.Context, serviceID, vehicleCategoryID int64) (dbpg.GetPriceTierRow, error)
}

// ServiceWithTiers bundles a service with its price tiers.
type ServiceWithTiers struct {
	Service dbpg.Service
	Tiers   []dbpg.ListPriceTiersByServiceRow
}

type CatalogueService struct {
	repo  CatalogueRepository
	authz *AuthzSvc
}

func NewCatalogueService(repo CatalogueRepository, authz *AuthzSvc) *CatalogueService {
	return &CatalogueService{
		repo:  repo,
		authz: authz,
	}
}

// ListCategories returns all service categories (public).
func (s *CatalogueService) ListCategories(ctx context.Context) ([]dbpg.ServiceCategory, error) {
	cats, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list categories", err)
	}
	return cats, nil
}

// ListAllServices returns all services (active and inactive) with their price tiers (admin only).
func (s *CatalogueService) ListAllServices(ctx context.Context, userID int64) ([]ServiceWithTiers, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, problems.New(problems.Unauthorized, "admin access required")
	}

	svcs, err := s.repo.ListAllServices(ctx)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list services", err)
	}

	result := make([]ServiceWithTiers, len(svcs))
	for i, svc := range svcs {
		tiers, err := s.repo.ListPriceTiersByService(ctx, svc.ID)
		if err != nil {
			return nil, problems.New(problems.Database, "failed to list price tiers", err)
		}
		result[i] = ServiceWithTiers{Service: svc, Tiers: tiers}
	}
	return result, nil
}

// ListServices returns all active services with their price tiers (public).
func (s *CatalogueService) ListServices(ctx context.Context) ([]ServiceWithTiers, error) {
	svcs, err := s.repo.ListServices(ctx)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list services", err)
	}

	result := make([]ServiceWithTiers, len(svcs))
	for i, svc := range svcs {
		tiers, err := s.repo.ListPriceTiersByService(ctx, svc.ID)
		if err != nil {
			return nil, problems.New(problems.Database, "failed to list price tiers", err)
		}
		result[i] = ServiceWithTiers{Service: svc, Tiers: tiers}
	}
	return result, nil
}

// GetServiceBySlug returns a service by slug with category name, options, and price tiers (public).
func (s *CatalogueService) GetServiceBySlug(ctx context.Context, slug string) (dbpg.GetServiceBySlugRow, []dbpg.ServiceOption, []dbpg.ListPriceTiersByServiceRow, error) {
	svc, err := s.repo.GetServiceBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return dbpg.GetServiceBySlugRow{}, nil, nil, problems.New(problems.NotExist, "service not found")
		}
		return dbpg.GetServiceBySlugRow{}, nil, nil, problems.New(problems.Database, "failed to get service", err)
	}

	opts, err := s.repo.ListServiceOptions(ctx, svc.ID)
	if err != nil {
		return dbpg.GetServiceBySlugRow{}, nil, nil, problems.New(problems.Database, "failed to list service options", err)
	}

	tiers, err := s.repo.ListPriceTiersByService(ctx, svc.ID)
	if err != nil {
		return dbpg.GetServiceBySlugRow{}, nil, nil, problems.New(problems.Database, "failed to list price tiers", err)
	}

	return svc, opts, tiers, nil
}

// CreateService creates a new service (admin only).
func (s *CatalogueService) CreateService(ctx context.Context, userID int64, params dbpg.CreateServiceParams) (dbpg.Service, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return dbpg.Service{}, err
	}
	if !isAdmin {
		return dbpg.Service{}, problems.New(problems.Unauthorized, "admin access required")
	}

	svc, err := s.repo.CreateService(ctx, params)
	if err != nil {
		return dbpg.Service{}, problems.New(problems.Database, "failed to create service", err)
	}
	return svc, nil
}

// UpdateService updates a service (admin only).
func (s *CatalogueService) UpdateService(ctx context.Context, userID int64, params dbpg.UpdateServiceParams) (dbpg.Service, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return dbpg.Service{}, err
	}
	if !isAdmin {
		return dbpg.Service{}, problems.New(problems.Unauthorized, "admin access required")
	}

	svc, err := s.repo.UpdateService(ctx, params)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return dbpg.Service{}, problems.New(problems.NotExist, "service not found")
		}
		return dbpg.Service{}, problems.New(problems.Database, "failed to update service", err)
	}
	return svc, nil
}

// DeleteService soft-deletes a service (admin only).
func (s *CatalogueService) DeleteService(ctx context.Context, userID int64, id int64) error {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return problems.New(problems.Unauthorized, "admin access required")
	}

	_, err = s.repo.DeleteService(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return problems.New(problems.NotExist, "service not found")
		}
		return problems.New(problems.Database, "failed to delete service", err)
	}
	return nil
}

// AddServiceOption adds an option to a service (admin only).
func (s *CatalogueService) AddServiceOption(ctx context.Context, userID int64, params dbpg.CreateServiceOptionParams) (dbpg.ServiceOption, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return dbpg.ServiceOption{}, err
	}
	if !isAdmin {
		return dbpg.ServiceOption{}, problems.New(problems.Unauthorized, "admin access required")
	}

	opt, err := s.repo.CreateServiceOption(ctx, params)
	if err != nil {
		return dbpg.ServiceOption{}, problems.New(problems.Database, "failed to create service option", err)
	}
	return opt, nil
}

// ListVehicleCategories returns all vehicle categories (public).
func (s *CatalogueService) ListVehicleCategories(ctx context.Context) ([]dbpg.VehicleCategory, error) {
	cats, err := s.repo.ListVehicleCategories(ctx)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list vehicle categories", err)
	}
	return cats, nil
}

// CreateVehicleCategory creates a vehicle category (admin only).
func (s *CatalogueService) CreateVehicleCategory(ctx context.Context, userID int64, params dbpg.CreateVehicleCategoryParams) (dbpg.VehicleCategory, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return dbpg.VehicleCategory{}, err
	}
	if !isAdmin {
		return dbpg.VehicleCategory{}, problems.New(problems.Unauthorized, "admin access required")
	}

	vc, err := s.repo.CreateVehicleCategory(ctx, params)
	if err != nil {
		return dbpg.VehicleCategory{}, problems.New(problems.Database, "failed to create vehicle category", err)
	}
	return vc, nil
}

// UpdateVehicleCategory updates a vehicle category (admin only).
func (s *CatalogueService) UpdateVehicleCategory(ctx context.Context, userID int64, params dbpg.UpdateVehicleCategoryParams) (dbpg.VehicleCategory, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return dbpg.VehicleCategory{}, err
	}
	if !isAdmin {
		return dbpg.VehicleCategory{}, problems.New(problems.Unauthorized, "admin access required")
	}

	vc, err := s.repo.UpdateVehicleCategory(ctx, params)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return dbpg.VehicleCategory{}, problems.New(problems.NotExist, "vehicle category not found")
		}
		return dbpg.VehicleCategory{}, problems.New(problems.Database, "failed to update vehicle category", err)
	}
	return vc, nil
}

// DeleteVehicleCategory deletes a vehicle category (admin only).
func (s *CatalogueService) DeleteVehicleCategory(ctx context.Context, userID int64, id int64) error {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return problems.New(problems.Unauthorized, "admin access required")
	}

	return s.repo.DeleteVehicleCategory(ctx, id)
}

// SetServicePriceTiers replaces all price tiers for a service (admin only).
func (s *CatalogueService) SetServicePriceTiers(ctx context.Context, userID int64, serviceID int64, tiers []dbpg.UpsertPriceTierParams) ([]dbpg.ListPriceTiersByServiceRow, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, problems.New(problems.Unauthorized, "admin access required")
	}

	// Delete existing tiers then upsert new ones
	if err := s.repo.DeletePriceTiersByService(ctx, serviceID); err != nil {
		return nil, problems.New(problems.Database, "failed to clear price tiers", err)
	}

	for _, t := range tiers {
		t.ServiceID = serviceID
		if _, err := s.repo.UpsertPriceTier(ctx, t); err != nil {
			return nil, problems.New(problems.Database, "failed to upsert price tier", err)
		}
	}

	result, err := s.repo.ListPriceTiersByService(ctx, serviceID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list price tiers", err)
	}
	return result, nil
}
