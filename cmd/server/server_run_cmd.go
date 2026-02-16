package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jpillora/backoff"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	ac "github.com/typewriterco/p402/internal/accesscontrol"
	"github.com/typewriterco/p402/internal/dbpg"
	fastmail "github.com/typewriterco/p402/internal/email/genericsmtp"
	gw "github.com/typewriterco/p402/internal/gateway/p402/v1"
	grpcsvr "github.com/typewriterco/p402/internal/grpc"
	"github.com/typewriterco/p402/internal/health"
	notification "github.com/typewriterco/p402/internal/notifications"
	pb "github.com/typewriterco/p402/internal/pb/p402/v1"
	"github.com/typewriterco/p402/internal/repos"
	"github.com/typewriterco/p402/internal/riverqueue"
	"github.com/typewriterco/p402/internal/services"
	"github.com/typewriterco/p402/internal/settings"
	"github.com/typewriterco/p402/internal/templater"
	thttp "github.com/typewriterco/p402/internal/transport/http"
	"github.com/typewriterco/p402/internal/workers"
	"github.com/urfave/cli/v2"
)

const (
	SERVER_DB_SCHEMA_NAME = "p402"
	RIVER_DB_SCHEMA_NAME  = "river"
	FGA_DB_SCHEMA_NAME    = "fga"
)

// getVersion returns version information from Go's built-in VCS data
// Automatically includes git commit, tag, and dirty status when built with Go 1.18+
func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown (no build info)"
	}

	version := "dev"
	revision := "unknown"
	modified := false
	buildTime := "unknown"

	// Extract VCS information from build settings
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
			if len(revision) > 7 {
				revision = revision[:7] // short hash
			}
		case "vcs.time":
			buildTime = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}

	// Try to find a version tag
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = info.Main.Version
	}

	// Build version string
	result := fmt.Sprintf("p402 %s (commit: %s, built: %s", version, revision, buildTime)
	if modified {
		result += ", dirty"
	}
	result += ")"

	return result
}

func serverRun(ctx *cli.Context) error {
	config := loadConfigFromCLI(ctx)

	setBaseLogger(ctx)

	log.Info().Str("opserver", "server init").Msg("")
	var dbStore *dbpg.Store
	var err error
	var dbCon *pgxpool.Pool
	b := &backoff.Backoff{Max: 5 * time.Minute}
	for {
		dbCon, err = dbpg.NewConnection(config.Database.ConnectionStringWithSchema(SERVER_DB_SCHEMA_NAME), ctx.App.Version)
		if err == nil {
			break
		}
		d := b.Duration()
		log.Warn().Err(err).Dur("retry_in", d).Msg("database connection failed, retrying")
		time.Sleep(d)
		continue
	}
	b.Reset()

	dbStore, err = dbpg.NewStore(dbCon)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database store")
	}
	log.Info().Int64("db_schema_version", dbStore.Version).Bool("dirty", dbStore.Dirty).Msg("Current DB Version")

	err = dbStore.CheckDB(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("database health check failed")
	}

	ds := dbpg.NewDataStore(dbStore)
	queries := dbpg.New(dbCon)

	fgaDBCon, err := dbpg.NewConnection(config.Database.ConnectionStringWithSchema(FGA_DB_SCHEMA_NAME), FGA_DB_SCHEMA_NAME)
	defer fgaDBCon.Close()

	if err != nil {
		log.Error().Err(err).Msg("failed to create a fga db client")
		return err
	}

	// Initialize new hierarchical settings service
	settingsService := settings.NewService(queries, log.Logger)

	// Check if development mode is enabled
	devMode := settings.NewDevMode(settingsService)
	if devMode.IsEnabled(context.Background()) {
		log.Warn().Msg("DEVELOPMENT MODE IS ENABLED - DO NOT USE IN PRODUCTION")

		// Log which dev mode features are enabled
		if devMode.SkipEmailVerification(context.Background()) {
			log.Warn().Str("feature", "skip_email_verification").Msg("dev mode feature enabled")
		}
		if devMode.DisableRateLimits(context.Background()) {
			log.Warn().Str("feature", "disable_rate_limits").Msg("dev mode feature enabled")
		}
		if devMode.AllowInsecureAuth(context.Background()) {
			log.Warn().Str("feature", "allow_insecure_auth").Msg("dev mode feature enabled")
		}
	}

	acClient, err := ac.New(context.Background(), fgaDBCon, log.Logger, settingsService)

	authzClient := services.NewAuthz(*acClient)

	if err != nil {
		log.Error().Err(err).Msg("cannot")
		os.Exit(1)
	}

	dr := repos.NewUserRepo(ds)

	authNService := services.NewAuthN(ds)

	userSvc, err := services.NewUserService(dr, authzClient, authNService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create UserService")
	}
	tpler, err := templater.NewTemplateManager(context.Background(), *dbStore)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create template manager")
	}

	riverDBCon, err := dbpg.NewConnection(config.Database.ConnectionStringWithSchema(RIVER_DB_SCHEMA_NAME), ctx.App.Version)

	ll := log.With().Logger()
	rq := riverqueue.New(riverDBCon, ll)

	smtpClient := fastmail.NewClient(settingsService)

	emailWorker := workers.NewEmailWorker(nil, smtpClient)

	emailWkrConfig := riverqueue.WorkerConfig{
		Name:       "email",
		Queue:      "email",
		MaxWorkers: 2,
	}

	riverqueue.Register(rq, emailWkrConfig, emailWorker)

	// Session cleanup worker
	sessionCleanupWorker := workers.NewSessionCleanupWorker(dbStore)
	maintenanceWkrConfig := riverqueue.WorkerConfig{
		Name:       "session_cleanup",
		Queue:      "maintenance",
		MaxWorkers: 1,
	}
	riverqueue.Register(rq, maintenanceWkrConfig, sessionCleanupWorker)

	err = rq.Start(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start river queuing")
	}

	// Schedule initial cleanup job
	_, err = rq.Client().InsertMany(context.Background(), []river.InsertManyParams{
		{Args: workers.SessionCleanupArgs{}},
	})
	if err != nil {
		log.Warn().Err(err).Msg("failed to schedule session cleanup job")
	}

	n := notification.NewNotifier(rq, tpler, config.DefaultFromEmail)

	signUpSvc := services.NewSignUp(userSvc, authNService, authzClient, settingsService)
	signUpSvc.Notifier = n

	// Create auth middleware for protecting endpoints
	authMiddleware := thttp.NewAuthMiddleware(authNService, authzClient)

	// Create health service with database checker
	dbChecker := health.NewDatabaseChecker(dbStore)
	healthSvc := health.NewService(dbChecker)

	// HTTP handlers (only for non-gRPC functionality)
	handlers := thttp.NewHandlers()
	handlers.SMTP = smtpClient // SMTP admin - not in proto

	// ========================================
	// gRPC Server Setup
	// ========================================

	// Create gRPC server with auth interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcsvr.AuthInterceptor(authNService)),
	)

	// Register gRPC services
	userGrpcSvc := grpcsvr.NewUserServiceServer(userSvc)
	pb.RegisterUserServiceServer(grpcServer, userGrpcSvc)

	authGrpcSvc := grpcsvr.NewAuthServiceServer(authNService, signUpSvc, config.BaseURL)
	pb.RegisterAuthServiceServer(grpcServer, authGrpcSvc)

	settingsRepo := repos.NewSettingsRepo(ds)
	settingsGrpcSvc := grpcsvr.NewSettingsServiceServer(settingsService, settingsRepo)
	pb.RegisterSettingsServiceServer(grpcServer, settingsGrpcSvc)

	// Enable gRPC reflection for grpcurl/grpcui
	reflection.Register(grpcServer)

	// Start gRPC server in background
	grpcAddr := fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create gRPC listener")
	}

	go func() {
		log.Info().Str("address", grpcAddr).Msg("gRPC server starting")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// ========================================
	// gRPC-Gateway HTTP Proxy Setup
	// ========================================

	gwCtx := context.Background()
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register gateway handlers - connect to gRPC server
	grpcEndpoint := fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port)

	err = gw.RegisterUserServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register UserService gateway")
	}

	err = gw.RegisterAuthServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register AuthService gateway")
	}

	err = gw.RegisterSettingsServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register SettingsService gateway")
	}

	// ========================================
	// HTTP Server with Gateway + Chi
	// ========================================

	server := thttp.NewServerWithGateway(config, healthSvc, handlers, authMiddleware, gwmux)
	err = server.Serve()

	if err != nil {
		return err
	}

	// Graceful shutdown of gRPC server
	grpcServer.GracefulStop()
	log.Info().Msg("gRPC server stopped gracefully")

	return nil
}
