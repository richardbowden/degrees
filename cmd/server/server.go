package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/jpillora/backoff"
	"github.com/typewriterco/p402/internal/config"
	"github.com/typewriterco/p402/internal/problems"

	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	apihttp "github.com/typewriterco/p402/internal/api/http"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/repos"
	"github.com/typewriterco/p402/internal/services"
)

type server struct {
	config     config.Config
	wg         sync.WaitGroup
	httpServer *http.Server

	startTime time.Time
}

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second

	appName = "p402-backend"
)

func newServer(c config.Config) (*server, error) {
	s := server{startTime: time.Now().UTC()}
	err := s.init(c)
	return &s, err
}

func (s *server) init(config config.Config) error {
	s.config = config
	log.Info().Str("opserver", "server init").Msg("")
	var dbStore *dbpg.Store
	var err error
	b := &backoff.Backoff{
		Max: 5 * time.Minute,
	}
	version := "p402 0.0.1-alpha"
	//todo(rich): does this need to be here or else where, also look at database retry
	for {
		dbStore, err = dbpg.NewStoreCreateCon(s.config.Database.ConnectionString(), version)
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

	//TODO(rich): services creation needs looking at.
	userSvc, err := services.NewUserService(dr)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create UserService")
	}

	//TODO(rich): handler creation needs looking at.
	uh := apihttp.NewUserHandler(userSvc)

	//*** Setup HTTP Server
	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(httplog.RequestLogger(log.Logger))
	mux.Use(middleware.AllowContentType("application/json"))
	api := humachi.New(mux, huma.DefaultConfig("p402", "0.0.0"))

	//huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
	//	details := make([]string, len(errs))
	//
	//	for i, err := range errs {
	//		details[i] = err.Error()
	//	}
	//	return &problems.OOO{Status: status, Detail: message, Details: details}
	//}

	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		p := problems.Problem{
			Status: status,
			Detail: msg,
		}
		for _, e := range errs {
			p.AddDetail(e)
		}
		return p
	}

	huma.AutoRegister(api, uh)

	apihttp.PrintRoutes(api)

	// HTTP Server Setup
	addr := fmt.Sprintf(":%d", s.config.HTTP.Port)
	log.Info().Str("address", addr).Msg("Server listening On")

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      mux,
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
	return nil
}

func (s *server) serveHttp() error {

	defer func() {
		log.Info().Msgf("server ran for %v", time.Since(s.startTime))
	}()

	shutdownErrorChan := make(chan error)
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- s.httpServer.Shutdown(ctx)
	}()

	log.Info().Msgf("start up took %v to reach running state", time.Since(s.startTime))

	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	log.Info().Msg("server stopped")

	s.wg.Wait()

	return nil
}
