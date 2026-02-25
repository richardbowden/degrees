package services

import (
	"context"
	"errors"
	"time"

	"github.com/richardbowden/degrees/internal/problems"
)

type CustomerProfile struct {
	ID        int64
	UserID    int64
	Phone     string
	Address   string
	Suburb    string
	Postcode  string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Vehicle struct {
	ID                int64
	CustomerID        int64
	Make              string
	Model             string
	Year              int32
	Colour            string
	Rego              string
	PaintType         string
	ConditionNotes    string
	IsPrimary         bool
	VehicleCategoryID int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CustomerRepository interface {
	CreateProfile(ctx context.Context, userID int64, phone, address, suburb, postcode, notes string) (CustomerProfile, error)
	GetProfileByUserID(ctx context.Context, userID int64) (CustomerProfile, error)
	UpdateProfile(ctx context.Context, id int64, phone, address, suburb, postcode, notes string) (CustomerProfile, error)
	ListCustomers(ctx context.Context, limit, offset int32) ([]CustomerProfile, error)
	CreateVehicle(ctx context.Context, customerID int64, make, model string, year int32, colour, rego, paintType, conditionNotes string, isPrimary bool, vehicleCategoryID int64) (Vehicle, error)
	GetVehicleByID(ctx context.Context, vehicleID int64) (Vehicle, error)
	ListVehiclesByCustomer(ctx context.Context, customerID int64) ([]Vehicle, error)
	UpdateVehicle(ctx context.Context, id int64, make, model string, year int32, colour, rego, paintType, conditionNotes string, isPrimary bool, vehicleCategoryID int64) (Vehicle, error)
	DeleteVehicle(ctx context.Context, vehicleID int64) error
}

type CustomerService struct {
	repo  CustomerRepository
	authz *AuthzSvc
}

func NewCustomerService(repo CustomerRepository, authz *AuthzSvc) *CustomerService {
	return &CustomerService{
		repo:  repo,
		authz: authz,
	}
}

// GetOrCreateProfile returns the customer profile for the given user, creating one if it doesn't exist.
func (s *CustomerService) GetOrCreateProfile(ctx context.Context, userID int64) (*CustomerProfile, error) {
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			// Auto-create profile on first access
			profile, err = s.repo.CreateProfile(ctx, userID, "", "", "", "", "")
			if err != nil {
				return nil, problems.New(problems.Database, "failed to create customer profile", err)
			}
			return &profile, nil
		}
		return nil, problems.New(problems.Database, "failed to get customer profile", err)
	}
	return &profile, nil
}

// UpdateProfile updates the authenticated user's customer profile.
func (s *CustomerService) UpdateProfile(ctx context.Context, userID int64, phone, address, suburb, postcode string) (*CustomerProfile, error) {
	// Get or create the profile first
	existing, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.repo.UpdateProfile(ctx, existing.ID, phone, address, suburb, postcode, existing.Notes)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "customer profile not found")
		}
		return nil, problems.New(problems.Database, "failed to update customer profile", err)
	}
	return &profile, nil
}

// AddVehicle adds a vehicle to the authenticated customer's profile.
func (s *CustomerService) AddVehicle(ctx context.Context, userID int64, make, model string, year int32, colour, rego, paintType, conditionNotes string, isPrimary bool, vehicleCategoryID int64) (*Vehicle, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	vehicle, err := s.repo.CreateVehicle(ctx, profile.ID, make, model, year, colour, rego, paintType, conditionNotes, isPrimary, vehicleCategoryID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to create vehicle", err)
	}
	return &vehicle, nil
}

// ListVehicles lists all vehicles for the authenticated customer.
func (s *CustomerService) ListVehicles(ctx context.Context, userID int64) ([]Vehicle, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	vehicles, err := s.repo.ListVehiclesByCustomer(ctx, profile.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list vehicles", err)
	}
	return vehicles, nil
}

// UpdateVehicle updates a vehicle, ensuring it belongs to the authenticated customer.
func (s *CustomerService) UpdateVehicle(ctx context.Context, userID, vehicleID int64, make, model string, year int32, colour, rego, paintType, conditionNotes string, isPrimary bool, vehicleCategoryID int64) (*Vehicle, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Verify vehicle belongs to this customer
	existing, err := s.repo.GetVehicleByID(ctx, vehicleID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "vehicle not found")
		}
		return nil, problems.New(problems.Database, "failed to get vehicle", err)
	}
	if existing.CustomerID != profile.ID {
		return nil, problems.New(problems.Unauthorized, "vehicle does not belong to this customer")
	}

	vehicle, err := s.repo.UpdateVehicle(ctx, vehicleID, make, model, year, colour, rego, paintType, conditionNotes, isPrimary, vehicleCategoryID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "vehicle not found")
		}
		return nil, problems.New(problems.Database, "failed to update vehicle", err)
	}
	return &vehicle, nil
}

// DeleteVehicle deletes a vehicle, ensuring it belongs to the authenticated customer.
func (s *CustomerService) DeleteVehicle(ctx context.Context, userID, vehicleID int64) error {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return err
	}

	// Verify vehicle belongs to this customer
	existing, err := s.repo.GetVehicleByID(ctx, vehicleID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return problems.New(problems.NotExist, "vehicle not found")
		}
		return problems.New(problems.Database, "failed to get vehicle", err)
	}
	if existing.CustomerID != profile.ID {
		return problems.New(problems.Unauthorized, "vehicle does not belong to this customer")
	}

	return s.repo.DeleteVehicle(ctx, vehicleID)
}

// ListCustomers lists all customer profiles (admin only).
func (s *CustomerService) ListCustomers(ctx context.Context, userID int64, limit, offset int32) ([]CustomerProfile, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, problems.New(problems.Internal, "failed to check admin permission", err)
	}
	if !isAdmin {
		return nil, problems.New(problems.Unauthorized, "admin access required")
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.ListCustomers(ctx, limit, offset)
}

// GetCustomer retrieves a customer profile with vehicles (admin only).
func (s *CustomerService) GetCustomer(ctx context.Context, userID, customerID int64) (*CustomerProfile, []Vehicle, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, nil, problems.New(problems.Internal, "failed to check admin permission", err)
	}
	if !isAdmin {
		return nil, nil, problems.New(problems.Unauthorized, "admin access required")
	}

	profile, err := s.repo.GetProfileByUserID(ctx, customerID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, nil, problems.New(problems.NotExist, "customer not found")
		}
		return nil, nil, problems.New(problems.Database, "failed to get customer", err)
	}

	vehicles, err := s.repo.ListVehiclesByCustomer(ctx, profile.ID)
	if err != nil {
		return nil, nil, problems.New(problems.Database, "failed to list customer vehicles", err)
	}

	return &profile, vehicles, nil
}
