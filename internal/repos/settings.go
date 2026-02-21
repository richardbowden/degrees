package repos

import (
	"context"

	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/settings"
)

type Settings struct {
	store dbpg.Storer
}

func NewSettingsRepo(store dbpg.Storer) *Settings {
	return &Settings{store: store}
}

func (r *Settings) ListAll(ctx context.Context) ([]settings.Setting, error) {
	dbSettings, err := r.store.ListAllSettings(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]settings.Setting, len(dbSettings))
	for i, dbSetting := range dbSettings {
		result[i] = convertDBSettingToDomain(dbSetting)
	}

	return result, nil
}

func convertDBSettingToDomain(db dbpg.Setting) settings.Setting {
	s := settings.Setting{
		ID:        db.ID,
		Scope:     db.Scope,
		Subsystem: db.Subsystem,
		Key:       db.Key,
		Value:     db.Value,
		CreatedAt: db.CreatedAt.Time,
		UpdatedAt: db.UpdatedAt.Time,
	}

	if db.OrganizationID.Valid {
		s.OrganizationID = &db.OrganizationID.Int64
	}
	if db.ProjectID.Valid {
		s.ProjectID = &db.ProjectID.Int64
	}
	if db.UserID.Valid {
		s.UserID = &db.UserID.Int64
	}
	if db.Description.Valid {
		s.Description = db.Description.String
	}
	if db.UpdatedBy.Valid {
		s.UpdatedBy = &db.UpdatedBy.Int64
	}

	return s
}
