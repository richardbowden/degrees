package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	openFGAMigrate "github.com/openfga/openfga/pkg/storage/migrate"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/richardbowden/degrees/internal/dbpg"
	migrator "github.com/richardbowden/degrees/internal/migrations"
	"github.com/richardbowden/degrees/sql/schema"
	"github.com/urfave/cli/v2"
)

func ensureSchemaExists(ctx context.Context, connString, schemaName string) error {
	tempConn, err := dbpg.NewConnection(connString, "")
	if err != nil {
		return fmt.Errorf("failed to create temp connection: %w", err)
	}
	defer tempConn.Close()

	_, err = tempConn.Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", pgx.Identifier{schemaName}.Sanitize()))
	if err != nil {
		return fmt.Errorf("failed to create schema %s: %w", schemaName, err)
	}
	return nil
}

func dbMigrate(ctx *cli.Context) error {
	dbCon := loadDBConfigFromCLI(ctx)

	schemas := []string{SERVER_DB_SCHEMA_NAME, FGA_DB_SCHEMA_NAME, RIVER_DB_SCHEMA_NAME}
	migrationCtx := context.Background()
	for _, schema := range schemas {
		if err := ensureSchemaExists(migrationCtx, dbCon.ConnectionString(), schema); err != nil {
			return fmt.Errorf("failed to create required schema: %w", err)
		}
	}

	mm, err := migrator.NewMigrator(schema.SQLMigrationFiles, dbCon.ConnectionStringWithSchema(SERVER_DB_SCHEMA_NAME))

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

	// Run River migrations
	riverDBCon, err := dbpg.NewConnection(dbCon.ConnectionStringWithSchema(RIVER_DB_SCHEMA_NAME), RIVER_DB_SCHEMA_NAME)
	if err != nil {
		return fmt.Errorf("failed to create river db connection: %w", err)
	}
	defer riverDBCon.Close()

	riverMigrator, err := rivermigrate.New(riverpgxv5.New(riverDBCon), &rivermigrate.Config{
		Schema: RIVER_DB_SCHEMA_NAME,
	})
	if err != nil {
		return fmt.Errorf("failed to create river migrator: %w", err)
	}

	_, err = riverMigrator.Migrate(context.Background(), rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{})
	if err != nil {
		return fmt.Errorf("failed to migrate river: %w", err)
	}

	return nil
}

func dbCurrentVersion(ctx *cli.Context) error {

	dbCon := loadDBConfigFromCLI(ctx)
	mm, err := migrator.NewMigrator(schema.SQLMigrationFiles, dbCon.ConnectionStringWithSchema(SERVER_DB_SCHEMA_NAME))

	if err != nil {
		return err
	}

	dbVersion, dbDirty, err := mm.Version()

	fmt.Printf("db_migration_version:%d, dirty: %v\n", dbVersion, dbDirty)

	return err
}
