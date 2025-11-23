package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jpillora/backoff"
	"github.com/rs/zerolog/log"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/fga"
	"github.com/typewriterco/p402/internal/repos"
	"github.com/typewriterco/p402/internal/services"
	"github.com/typewriterco/p402/internal/settings"
	thttp "github.com/typewriterco/p402/internal/transport/http"
	"github.com/urfave/cli/v2"
)

const (
	FGA_DB_SCHEMA_NAME = "fga"
)

var version = "p402 0.0.1-alpha"

func serverRun(ctx *cli.Context) error {
	config := loadConfigFromCLI(ctx)

	setBaseLogger(ctx)

	log.Info().Str("opserver", "server init").Msg("")
	var dbStore *dbpg.Store
	var err error
	var dbCon *pgxpool.Pool
	b := &backoff.Backoff{Max: 5 * time.Minute}
	for {
		dbCon, err = dbpg.NewConnection(config.Database.ConnectionString(), ctx.App.Version)
		if err == nil {
			break
		}
		d := b.Duration()
		fmt.Printf("%s, reconnecting in %s\n", err, d)
		time.Sleep(d)
		continue
	}
	b.Reset()

	dbStore = dbpg.NewStore(dbCon)
	log.Info().Int64("db_schema_version", dbStore.Version).Bool("dirty", dbStore.Dirty).Msg("Current DB Version")

	dbStore.CheckDB(context.Background())

	ds := dbpg.NewDataStore(dbStore)

	fgaDBCon, err := dbpg.NewConnection(config.Database.ConnectionStringWithSchema(FGA_DB_SCHEMA_NAME), FGA_DB_SCHEMA_NAME)
	defer fgaDBCon.Close()

	if err != nil {
		log.Error().Err(err).Msg("failed to create a fga db client")
		return err
	}

	settings := settings.New(dbStore)

	fgaClient, err := fga.New(context.Background(), fgaDBCon, log.Logger, settings)

	if err != nil {
		log.Error().Err(err).Msg("cannot")
		os.Exit(1)
	}

	fgaClient.ListFiles()

	// Set dev overrides to help run things locally
	// services.DevSkipUserVerification = s.config.devOverrideConfig.SkipUserConfirm

	dr := repos.NewAccountsRepo(ds)

	userSvc, err := services.NewUserService(dr)
	userHandlers := thttp.NewUserHandler(userSvc)

	authSvc := services.NewAuth(dr)
	authHandlers := thttp.NewAuth(authSvc)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create UserService")
	}

	handlers := thttp.NewHandlers()
	handlers.Users = userHandlers
	handlers.Auth = authHandlers

	server := thttp.NewServer(config, handlers)
	err = server.Serve()

	if err != nil {
		return err
	}

	return nil
}
