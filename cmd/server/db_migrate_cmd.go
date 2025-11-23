package main

import (
	"fmt"
	"strings"

	openFGAMigrate "github.com/openfga/openfga/pkg/storage/migrate"

	migrator "github.com/typewriterco/p402/internal/migrations"
	"github.com/typewriterco/p402/sql/schema"
	"github.com/urfave/cli/v2"
)

func dbMigrate(ctx *cli.Context) error {
	dbCon := loadDBConfigFromCLI(ctx)
	mm, err := migrator.NewMigrator(schema.SQLMigrationFiles, dbCon.ConnectionString())

	if err != nil {
		return err
	}

	err = mm.Migrate()
	if err != nil {
		if !strings.Contains(err.Error(), "no change") {
			return err
		}
	}

	fgaConfig := openFGAMigrate.MigrationConfig{
		Engine:  "postgres",
		URI:     dbCon.ConnectionStringWithSchema(FGA_DB_SCHEMA_NAME),
		Verbose: false,
	}

	err = openFGAMigrate.RunMigrations(fgaConfig)

	if err != nil {
		return fmt.Errorf("failed to migrate fga %w", err)
	}

	return nil
}

func dbCurrentVersion(ctx *cli.Context) error {

	dbCon := loadDBConfigFromCLI(ctx)
	mm, err := migrator.NewMigrator(schema.SQLMigrationFiles, dbCon.ConnectionString())

	if err != nil {
		return err
	}

	dbVersion, dbDirty, err := mm.Version()

	fmt.Printf("db_migration_version:%d, dirty: %v\n", dbVersion, dbDirty)

	return err
}
