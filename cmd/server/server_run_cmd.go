package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jpillora/backoff"
	"github.com/rs/zerolog/log"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/repos"
	"github.com/typewriterco/p402/internal/services"
	thttp "github.com/typewriterco/p402/internal/transport/http"
	"github.com/urfave/cli/v2"
)

var version = "p402 0.0.1-alpha"

func serverRun(ctx *cli.Context) error {
	config := loadConfigFromCLI(ctx)

	setBaseLogger(ctx)

	log.Info().Str("opserver", "server init").Msg("")
	var dbStore *dbpg.Store
	var err error

	b := &backoff.Backoff{Max: 5 * time.Minute}
	for {
		dbStore, err = dbpg.NewStoreCreateCon(config.Database.ConnectionString(), version)
		if err == nil {
			break
		}
		d := b.Duration()
		fmt.Printf("%s, reconnecting in %s\n", err, d)
		time.Sleep(d)
		continue
	}
	b.Reset()

	log.Info().Int64("db_schema_version", dbStore.Version).Bool("dirty", dbStore.Dirty).Msg("Current DB Version")

	dbStore.CheckDB(context.Background())

	ds := dbpg.NewDataStore(dbStore)

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
