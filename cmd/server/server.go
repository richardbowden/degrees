package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	apihttp "github.com/typewriterco/p402/internal/api/http"
	"github.com/typewriterco/p402/internal/dbpg"
	migrator "github.com/typewriterco/p402/internal/migrations"
	"github.com/typewriterco/p402/internal/repos"
	"github.com/typewriterco/p402/internal/services"
	"github.com/typewriterco/p402/sql/schema"
)

type server struct {
	config         config
	wg             sync.WaitGroup
	httpServer     *http.Server
	router         *chi.Mux
	accountService *services.AccountService
	accountHandler *apihttp.AccountHandler

	debugHandler *apihttp.DebugHandler

	start_time time.Time

	pg *pgxpool.Pool //TODO(rich): this is temp, the db package should expose db stats
}

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second

	appName = "p402-backend"
)

func NewServer(c config) (*server, error) {
	s := server{start_time: time.Now().UTC()}
	err := s.init(c)
	return &s, err
}

func (s *server) init(config config) error {
	s.config = config
	log.Info().Str("opserver", "server init").Msg("")

	//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
	//This is tempoary.... need to add a better db version check on start up.
	//In a read only mode. We do not want to migrate as part of server startup
	//due to there will be more than one copy of the server running or starting up.
	//
	//maybe compare the embded sql files against the version of the DB
	m, err := migrator.NewMigrator(schema.SQLMigrationFiles, s.config.db.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create a Migrator")
	}

	version, dirty, err := m.Version()

	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	log.Info().Uint("db_schema_version", version).Bool("dirty", dirty).Msg("Current DB Version")

	if dirty {
		log.Fatal().Bool("dirty", dirty).Msg("database is dirty, server unable to start until database schema is fixed")
	}

	defer m.Close()
	//##############################################################

	dbStore, err := dbpg.NewStoreCreateCon(s.config.db.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create a new store")
	}

	s.pg = dbStore.SQLPool

	dbStore.CheckDB(context.Background())

	ds := dbpg.NewDataStore(dbStore)

	// Set dev overrides to help run things locally
	// services.DevSkipUserVerification = s.config.devOverrideConfig.SkipUserConfirm

	dr := repos.NewAccountsRepo(ds)

	//TODO(rich): services creation needs looking at.
	accountsSvc, err := services.NewAccountService(ds, dr)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create AccountService")
	}

	s.accountService = accountsSvc

	//TODO(rich): handler creation needs looking at.
	ah := apihttp.NewAccountHandler(accountsSvc)

	s.accountHandler = ah

	s.debugHandler = &apihttp.DebugHandler{}

	s.router = chi.NewRouter()
	s.router.Use(middleware.RequestID)
	s.router.Use(httplog.RequestLogger(log.Logger))
	s.router.Use(middleware.AllowContentType("application/json"))

	s.router.Mount(
		"/", s.Endpoints(),
	)

	s.router.Group(func(r chi.Router) {
		r.Use(apihttp.IsAuthed())
		r.Get("/profile", profile)
	})

	addr := fmt.Sprintf(":%s", s.config.httpPort)
	log.Info().Str("address", addr).Msg("Server listening On")

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
	return nil
}

func (s *server) serveHttp() error {

	defer func() {
		log.Info().Msgf("server ran for %v", time.Since(s.start_time))
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

	log.Info().Msgf("start up took %v to reach running state", time.Since(s.start_time))

	s.walkRoutes()

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

func profile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("this is a test response"))
}
