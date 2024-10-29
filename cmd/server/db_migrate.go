package main

import (
	"fmt"

	migrator "github.com/typewriterco/p402/internal/migrations"
	"github.com/typewriterco/p402/sql/schema"
	"github.com/urfave/cli/v2"
)

func db_migrate(ctx *cli.Context) error {
	db_con := DBConfigFromCTX(ctx)
	mm, err := migrator.NewMigrator(schema.SQLMigrationFiles, db_con.ConnectionString())

	if err != nil {
		return err
	}

	err = mm.Migrate()

	if err != nil {
		return err
	}

	return nil
}

func db_current_version(ctx *cli.Context) error {

	db_con := DBConfigFromCTX(ctx)
	mm, err := migrator.NewMigrator(schema.SQLMigrationFiles, db_con.ConnectionString())

	if err != nil {
		return err
	}

	version, dirty, err := mm.Version()

	fmt.Printf("db_migration_version:%d, dirty: %v\n", version, dirty)

	return err
}
