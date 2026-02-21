package repos

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/services"
)

type History struct {
	store dbpg.Storer
}

func NewHistoryRepo(store dbpg.Storer) *History {
	return &History{
		store: store,
	}
}

func (r *History) CreateServiceRecord(ctx context.Context, bookingID, customerID, vehicleID int64, completedDate time.Time) (services.ServiceRecord, error) {
	dbRecord, err := r.store.CreateServiceRecord(ctx, dbpg.CreateServiceRecordParams{
		BookingID:     bookingID,
		CustomerID:    customerID,
		VehicleID:     pgtype.Int8{Int64: vehicleID, Valid: vehicleID > 0},
		CompletedDate: pgtype.Timestamptz{Time: completedDate, Valid: !completedDate.IsZero()},
	})
	if err != nil {
		return services.ServiceRecord{}, err
	}
	return dbRecordToService(dbRecord), nil
}

func (r *History) GetServiceRecordByID(ctx context.Context, id int64) (services.ServiceRecord, error) {
	dbRecord, err := r.store.GetServiceRecordByID(ctx, dbpg.GetServiceRecordByIDParams{
		ID: id,
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.ServiceRecord{}, services.ErrNoRecord
		}
		return services.ServiceRecord{}, err
	}
	return dbRecordToService(dbRecord), nil
}

func (r *History) ListServiceRecordsByCustomer(ctx context.Context, customerID int64) ([]services.ServiceRecord, error) {
	dbRecords, err := r.store.ListServiceRecordsByCustomer(ctx, dbpg.ListServiceRecordsByCustomerParams{
		CustomerID: customerID,
	})
	if err != nil {
		return nil, err
	}
	records := make([]services.ServiceRecord, len(dbRecords))
	for i, rec := range dbRecords {
		records[i] = dbRecordToService(rec)
	}
	return records, nil
}

func (r *History) CreateServiceNote(ctx context.Context, serviceRecordID int64, noteType, content string, isVisibleToCustomer bool, createdBy int64) (services.ServiceNote, error) {
	dbNote, err := r.store.CreateServiceNote(ctx, dbpg.CreateServiceNoteParams{
		ServiceRecordID:     serviceRecordID,
		NoteType:            noteType,
		Content:             content,
		IsVisibleToCustomer: isVisibleToCustomer,
		CreatedBy:           pgtype.Int8{Int64: createdBy, Valid: createdBy > 0},
	})
	if err != nil {
		return services.ServiceNote{}, err
	}
	return dbNoteToService(dbNote), nil
}

func (r *History) ListServiceNotes(ctx context.Context, serviceRecordID int64) ([]services.ServiceNote, error) {
	dbNotes, err := r.store.ListServiceNotes(ctx, dbpg.ListServiceNotesParams{
		ServiceRecordID: serviceRecordID,
	})
	if err != nil {
		return nil, err
	}
	notes := make([]services.ServiceNote, len(dbNotes))
	for i, n := range dbNotes {
		notes[i] = dbNoteToService(n)
	}
	return notes, nil
}

func (r *History) CreateServiceProductUsed(ctx context.Context, serviceRecordID int64, productName, notes string) (services.ServiceProductUsed, error) {
	dbProduct, err := r.store.CreateServiceProductUsed(ctx, dbpg.CreateServiceProductUsedParams{
		ServiceRecordID: serviceRecordID,
		ProductName:     productName,
		Notes:           dbpg.StringToPGString(notes),
	})
	if err != nil {
		return services.ServiceProductUsed{}, err
	}
	return dbProductToService(dbProduct), nil
}

func (r *History) ListServiceProductsUsed(ctx context.Context, serviceRecordID int64) ([]services.ServiceProductUsed, error) {
	dbProducts, err := r.store.ListServiceProductsUsed(ctx, dbpg.ListServiceProductsUsedParams{
		ServiceRecordID: serviceRecordID,
	})
	if err != nil {
		return nil, err
	}
	products := make([]services.ServiceProductUsed, len(dbProducts))
	for i, p := range dbProducts {
		products[i] = dbProductToService(p)
	}
	return products, nil
}

func (r *History) ListServicePhotos(ctx context.Context, serviceRecordID int64) ([]services.ServicePhoto, error) {
	dbPhotos, err := r.store.ListServicePhotos(ctx, dbpg.ListServicePhotosParams{
		ServiceRecordID: serviceRecordID,
	})
	if err != nil {
		return nil, err
	}
	photos := make([]services.ServicePhoto, len(dbPhotos))
	for i, p := range dbPhotos {
		photos[i] = dbPhotoToService(p)
	}
	return photos, nil
}

func dbRecordToService(r dbpg.ServiceRecord) services.ServiceRecord {
	return services.ServiceRecord{
		ID:            r.ID,
		BookingID:     r.BookingID,
		CustomerID:    r.CustomerID,
		VehicleID:     r.VehicleID.Int64,
		CompletedDate: r.CompletedDate.Time,
		CreatedAt:     r.CreatedAt.Time,
		UpdatedAt:     r.UpdatedAt.Time,
	}
}

func dbNoteToService(n dbpg.ServiceNote) services.ServiceNote {
	return services.ServiceNote{
		ID:                  n.ID,
		ServiceRecordID:     n.ServiceRecordID,
		NoteType:            n.NoteType,
		Content:             n.Content,
		IsVisibleToCustomer: n.IsVisibleToCustomer,
		CreatedBy:           n.CreatedBy.Int64,
		CreatedAt:           n.CreatedAt.Time,
	}
}

func dbProductToService(p dbpg.ServiceProductsUsed) services.ServiceProductUsed {
	return services.ServiceProductUsed{
		ID:              p.ID,
		ServiceRecordID: p.ServiceRecordID,
		ProductName:     p.ProductName,
		Notes:           p.Notes.String,
	}
}

func dbPhotoToService(p dbpg.ServicePhoto) services.ServicePhoto {
	return services.ServicePhoto{
		ID:              p.ID,
		ServiceRecordID: p.ServiceRecordID,
		PhotoType:       p.PhotoType,
		URL:             p.Url,
		Caption:         p.Caption.String,
		CreatedAt:       p.CreatedAt.Time,
	}
}
