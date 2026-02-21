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
	ListServicesByCategory(ctx context.Context, categoryID int64) ([]dbpg.Service, error)
	GetServiceBySlug(ctx context.Context, slug string) (dbpg.GetServiceBySlugRow, error)
	GetServiceByID(ctx context.Context, id int64) (dbpg.Service, error)
	CreateService(ctx context.Context, params dbpg.CreateServiceParams) (dbpg.Service, error)
	UpdateService(ctx context.Context, params dbpg.UpdateServiceParams) (dbpg.Service, error)
	DeleteService(ctx context.Context, id int64) (dbpg.Service, error)
	ListServiceOptions(ctx context.Context, serviceID int64) ([]dbpg.ServiceOption, error)
	CreateServiceOption(ctx context.Context, params dbpg.CreateServiceOptionParams) (dbpg.ServiceOption, error)
	UpdateServiceOption(ctx context.Context, params dbpg.UpdateServiceOptionParams) (dbpg.ServiceOption, error)
	DeleteServiceOption(ctx context.Context, id int64) (dbpg.ServiceOption, error)
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

// ListServices returns all active services (public).
func (s *CatalogueService) ListServices(ctx context.Context) ([]dbpg.Service, error) {
	svcs, err := s.repo.ListServices(ctx)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list services", err)
	}
	return svcs, nil
}

// GetServiceBySlug returns a service by slug with category name (public).
func (s *CatalogueService) GetServiceBySlug(ctx context.Context, slug string) (dbpg.GetServiceBySlugRow, []dbpg.ServiceOption, error) {
	svc, err := s.repo.GetServiceBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return dbpg.GetServiceBySlugRow{}, nil, problems.New(problems.NotExist, "service not found")
		}
		return dbpg.GetServiceBySlugRow{}, nil, problems.New(problems.Database, "failed to get service", err)
	}

	opts, err := s.repo.ListServiceOptions(ctx, svc.ID)
	if err != nil {
		return dbpg.GetServiceBySlugRow{}, nil, problems.New(problems.Database, "failed to list service options", err)
	}

	return svc, opts, nil
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
