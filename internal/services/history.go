package services

import (
	"context"
	"errors"
	"time"

	"github.com/richardbowden/degrees/internal/problems"
)

type ServiceRecord struct {
	ID            int64
	BookingID     int64
	CustomerID    int64
	VehicleID     int64
	CompletedDate time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ServiceNote struct {
	ID                  int64
	ServiceRecordID     int64
	NoteType            string
	Content             string
	IsVisibleToCustomer bool
	CreatedBy           int64
	CreatedAt           time.Time
}

type ServiceProductUsed struct {
	ID              int64
	ServiceRecordID int64
	ProductName     string
	Notes           string
}

type ServicePhoto struct {
	ID              int64
	ServiceRecordID int64
	PhotoType       string
	URL             string
	Caption         string
	CreatedAt       time.Time
}

type ServiceRecordDetail struct {
	Record   ServiceRecord
	Notes    []ServiceNote
	Products []ServiceProductUsed
	Photos   []ServicePhoto
}

type HistoryRepository interface {
	CreateServiceRecord(ctx context.Context, bookingID, customerID, vehicleID int64, completedDate time.Time) (ServiceRecord, error)
	GetServiceRecordByID(ctx context.Context, id int64) (ServiceRecord, error)
	ListServiceRecordsByCustomer(ctx context.Context, customerID int64) ([]ServiceRecord, error)
	CreateServiceNote(ctx context.Context, serviceRecordID int64, noteType, content string, isVisibleToCustomer bool, createdBy int64) (ServiceNote, error)
	ListServiceNotes(ctx context.Context, serviceRecordID int64) ([]ServiceNote, error)
	CreateServiceProductUsed(ctx context.Context, serviceRecordID int64, productName, notes string) (ServiceProductUsed, error)
	ListServiceProductsUsed(ctx context.Context, serviceRecordID int64) ([]ServiceProductUsed, error)
	ListServicePhotos(ctx context.Context, serviceRecordID int64) ([]ServicePhoto, error)
}

type HistoryService struct {
	repo     HistoryRepository
	authz    *AuthzSvc
	customer CustomerRepository
}

func NewHistoryService(repo HistoryRepository, authz *AuthzSvc, customer CustomerRepository) *HistoryService {
	return &HistoryService{
		repo:     repo,
		authz:    authz,
		customer: customer,
	}
}

var validNoteTypes = map[string]bool{
	"condition":      true,
	"treatment":      true,
	"recommendation": true,
	"follow_up":      true,
}

// ListMyHistory returns service records for the authenticated user, with notes filtered to visible-only.
func (s *HistoryService) ListMyHistory(ctx context.Context, userID int64) ([]ServiceRecordDetail, error) {
	profile, err := s.customer.GetProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return []ServiceRecordDetail{}, nil
		}
		return nil, problems.New(problems.Database, "failed to get customer profile", err)
	}

	records, err := s.repo.ListServiceRecordsByCustomer(ctx, profile.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list service records", err)
	}

	details := make([]ServiceRecordDetail, len(records))
	for i, rec := range records {
		detail, err := s.loadRecordDetail(ctx, rec, false)
		if err != nil {
			return nil, err
		}
		details[i] = *detail
	}
	return details, nil
}

// GetServiceRecord returns a service record with details. Customers can only see their own records and visible notes.
func (s *HistoryService) GetServiceRecord(ctx context.Context, userID, recordID int64) (*ServiceRecordDetail, error) {
	record, err := s.repo.GetServiceRecordByID(ctx, recordID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "service record not found")
		}
		return nil, problems.New(problems.Database, "failed to get service record", err)
	}

	// Check if admin
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, problems.New(problems.Internal, "failed to check admin permission", err)
	}

	if !isAdmin {
		// For non-admins, verify they own this record
		profile, err := s.customer.GetProfileByUserID(ctx, userID)
		if err != nil {
			if errors.Is(err, ErrNoRecord) {
				return nil, problems.New(problems.Unauthorized, "access denied")
			}
			return nil, problems.New(problems.Database, "failed to get customer profile", err)
		}
		if record.CustomerID != profile.ID {
			return nil, problems.New(problems.Unauthorized, "access denied")
		}
	}

	return s.loadRecordDetail(ctx, record, isAdmin)
}

// ListCustomerHistory lists all service records for a customer (admin only), with all notes.
func (s *HistoryService) ListCustomerHistory(ctx context.Context, userID, customerID int64) ([]ServiceRecordDetail, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, problems.New(problems.Internal, "failed to check admin permission", err)
	}
	if !isAdmin {
		return nil, problems.New(problems.Unauthorized, "admin access required")
	}

	records, err := s.repo.ListServiceRecordsByCustomer(ctx, customerID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list service records", err)
	}

	details := make([]ServiceRecordDetail, len(records))
	for i, rec := range records {
		detail, err := s.loadRecordDetail(ctx, rec, true)
		if err != nil {
			return nil, err
		}
		details[i] = *detail
	}
	return details, nil
}

// CreateServiceRecord creates a new service record (admin only).
func (s *HistoryService) CreateServiceRecord(ctx context.Context, userID, bookingID, customerID, vehicleID int64, completedDate time.Time) (*ServiceRecord, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, problems.New(problems.Internal, "failed to check admin permission", err)
	}
	if !isAdmin {
		return nil, problems.New(problems.Unauthorized, "admin access required")
	}

	record, err := s.repo.CreateServiceRecord(ctx, bookingID, customerID, vehicleID, completedDate)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to create service record", err)
	}
	return &record, nil
}

// AddServiceNote adds a note to a service record (admin only).
func (s *HistoryService) AddServiceNote(ctx context.Context, userID, serviceRecordID int64, noteType, content string, isVisibleToCustomer bool) (*ServiceNote, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, problems.New(problems.Internal, "failed to check admin permission", err)
	}
	if !isAdmin {
		return nil, problems.New(problems.Unauthorized, "admin access required")
	}

	if !validNoteTypes[noteType] {
		return nil, problems.New(problems.Validation, "invalid note type, must be one of: condition, treatment, recommendation, follow_up")
	}

	// Verify service record exists
	_, err = s.repo.GetServiceRecordByID(ctx, serviceRecordID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "service record not found")
		}
		return nil, problems.New(problems.Database, "failed to get service record", err)
	}

	note, err := s.repo.CreateServiceNote(ctx, serviceRecordID, noteType, content, isVisibleToCustomer, userID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to create service note", err)
	}
	return &note, nil
}

// AddProductUsed adds a product used entry to a service record (admin only).
func (s *HistoryService) AddProductUsed(ctx context.Context, userID, serviceRecordID int64, productName, notes string) (*ServiceProductUsed, error) {
	isAdmin, err := s.authz.IsSystemAdmin(ctx, userID)
	if err != nil {
		return nil, problems.New(problems.Internal, "failed to check admin permission", err)
	}
	if !isAdmin {
		return nil, problems.New(problems.Unauthorized, "admin access required")
	}

	// Verify service record exists
	_, err = s.repo.GetServiceRecordByID(ctx, serviceRecordID)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "service record not found")
		}
		return nil, problems.New(problems.Database, "failed to get service record", err)
	}

	product, err := s.repo.CreateServiceProductUsed(ctx, serviceRecordID, productName, notes)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to add product used", err)
	}
	return &product, nil
}

// loadRecordDetail loads notes, products, and photos for a service record.
// If includeAllNotes is false, only customer-visible notes are included.
func (s *HistoryService) loadRecordDetail(ctx context.Context, record ServiceRecord, includeAllNotes bool) (*ServiceRecordDetail, error) {
	notes, err := s.repo.ListServiceNotes(ctx, record.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list service notes", err)
	}

	if !includeAllNotes {
		visibleNotes := make([]ServiceNote, 0, len(notes))
		for _, n := range notes {
			if n.IsVisibleToCustomer {
				visibleNotes = append(visibleNotes, n)
			}
		}
		notes = visibleNotes
	}

	products, err := s.repo.ListServiceProductsUsed(ctx, record.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list products used", err)
	}

	photos, err := s.repo.ListServicePhotos(ctx, record.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list service photos", err)
	}

	return &ServiceRecordDetail{
		Record:   record,
		Notes:    notes,
		Products: products,
		Photos:   photos,
	}, nil
}
