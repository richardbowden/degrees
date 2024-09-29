package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
)

type Migrator struct {
	fs embed.FS

	m *migrate.Migrate

	sqlCon *sql.DB
}

func (m *Migrator) Close() {
	m.m.Close()
}

func NewMigrator(fs embed.FS, dbCon string) (*Migrator, error) {

	m := &Migrator{fs: fs}

	files, err := iofs.New(fs, ".")

	if err != nil {
		return nil, err
	}

	mm, err := migrate.NewWithSourceInstance("iofs", files, dbCon)
	if err != nil {
		return nil, err
	}

	m.m = mm

	return m, nil
}

func (m *Migrator) Migrate() error {

	err := m.m.Up()

	if err != nil {
		return fmt.Errorf("failed to migrate db, %s", err)
	}

	return nil
}

func (m *Migrator) Version() (uint, bool, error) {

	version, dirty, err := m.m.Version()

	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, err
	}

	if errors.Is(err, migrate.ErrNilVersion) {
		log.Fatal().Msg("Database has not been migrated.")
	}
	return version, dirty, nil
}
