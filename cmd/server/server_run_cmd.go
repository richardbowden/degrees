package main

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jpillora/backoff"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	ac "github.com/richardbowden/degrees/internal/accesscontrol"
	"github.com/richardbowden/degrees/internal/dbpg"
	fastmail "github.com/richardbowden/degrees/internal/email/genericsmtp"
	gw "github.com/richardbowden/degrees/internal/gateway/degrees/v1"
	grpcsvr "github.com/richardbowden/degrees/internal/grpc"
	"github.com/richardbowden/degrees/internal/health"
	notification "github.com/richardbowden/degrees/internal/notifications"
	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/repos"
	"github.com/richardbowden/degrees/internal/riverqueue"
	"github.com/richardbowden/degrees/internal/services"
	"github.com/richardbowden/degrees/internal/settings"
	"github.com/richardbowden/degrees/internal/templater"
	thttp "github.com/richardbowden/degrees/internal/transport/http"
	"github.com/richardbowden/degrees/internal/workers"
	"github.com/urfave/cli/v2"
)

const (
	SERVER_DB_SCHEMA_NAME = "degrees"
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
	defer dbCon.Close()

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
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create fga db connection")
	}
	defer fgaDBCon.Close()

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
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create access control client")
	}

	authzClient := services.NewAuthz(*acClient)

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
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create river db connection")
	}

	ll := log.With().Logger()
	rq := riverqueue.New(riverDBCon, ll)

	smtpClient := fastmail.NewClient(settingsService)

	emailWorker := workers.NewEmailWorker(nil, smtpClient)

	emailWkrConfig := riverqueue.WorkerConfig{
		Name:       "email",
		Queue:      "email",
		MaxWorkers: 2,
	}

	if err := riverqueue.Register(rq, emailWkrConfig, emailWorker); err != nil {
		log.Fatal().Err(err).Msg("failed to register email worker")
	}

	// Session cleanup worker
	sessionCleanupWorker := workers.NewSessionCleanupWorker(dbStore)
	maintenanceWkrConfig := riverqueue.WorkerConfig{
		Name:       "session_cleanup",
		Queue:      "maintenance",
		MaxWorkers: 1,
	}
	if err := riverqueue.Register(rq, maintenanceWkrConfig, sessionCleanupWorker); err != nil {
		log.Fatal().Err(err).Msg("failed to register session cleanup worker")
	}

	// Schedule session cleanup to run hourly (and once on start)
	riverqueue.AddPeriodicJob(rq, 1*time.Hour, workers.SessionCleanupArgs{})

	n := notification.NewNotifier(rq, tpler, config.DefaultFromEmail)

	// Booking confirmation worker
	bookingConfirmationWorker := workers.NewBookingConfirmationWorker(n)
	bookingWkrConfig := riverqueue.WorkerConfig{
		Name:       "booking_confirmation",
		Queue:      "booking",
		MaxWorkers: 2,
	}
	if err := riverqueue.Register(rq, bookingWkrConfig, bookingConfirmationWorker); err != nil {
		log.Fatal().Err(err).Msg("failed to register booking confirmation worker")
	}

	err = rq.Start(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start river queuing")
	}

	signUpSvc := services.NewSignUp(userSvc, authNService, authzClient, settingsService)
	signUpSvc.Notifier = n

	// Create auth middleware for protecting endpoints
	authMiddleware := thttp.NewAuthMiddleware(authNService, authzClient)

	// Create health service with database checker
	dbChecker := health.NewDatabaseChecker(dbStore)
	healthSvc := health.NewService(dbChecker)

	// ========================================
	// gRPC Server Setup
	// ========================================

	// Create gRPC server with auth interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcsvr.AuthInterceptor(authNService)),
	)

	// Register gRPC services
	userGrpcSvc := grpcsvr.NewUserServiceServer(userSvc, authzClient)
	pb.RegisterUserServiceServer(grpcServer, userGrpcSvc)

	authGrpcSvc := grpcsvr.NewAuthServiceServer(authNService, signUpSvc, config.BaseURL)
	pb.RegisterAuthServiceServer(grpcServer, authGrpcSvc)

	settingsRepo := repos.NewSettingsRepo(ds)
	settingsGrpcSvc := grpcsvr.NewSettingsServiceServer(settingsService, settingsRepo, authzClient)
	pb.RegisterSettingsServiceServer(grpcServer, settingsGrpcSvc)

	smtpGrpcSvc := grpcsvr.NewSMTPServiceServer(smtpClient, authzClient)
	pb.RegisterSMTPServiceServer(grpcServer, smtpGrpcSvc)

	// Catalogue service
	catalogueRepo := repos.NewCatalogueRepo(ds)
	catalogueSvc := services.NewCatalogueService(catalogueRepo, authzClient)
	catalogueGrpcSvc := grpcsvr.NewCatalogueServiceServer(catalogueSvc)
	pb.RegisterCatalogueServiceServer(grpcServer, catalogueGrpcSvc)

	// Cart service
	cartRepo := repos.NewCartRepo(ds)
	cartSvc := services.NewCartService(cartRepo)
	cartGrpcSvc := grpcsvr.NewCartServiceServer(cartSvc)
	pb.RegisterCartServiceServer(grpcServer, cartGrpcSvc)

	// Customer service
	customerRepo := repos.NewCustomerRepo(ds)
	customerSvc := services.NewCustomerService(customerRepo, authzClient)
	customerGrpcSvc := grpcsvr.NewCustomerServiceServer(customerSvc)
	pb.RegisterCustomerServiceServer(grpcServer, customerGrpcSvc)

	// History service
	historyRepo := repos.NewHistoryRepo(ds)
	historySvc := services.NewHistoryService(historyRepo, authzClient, customerRepo)
	historyGrpcSvc := grpcsvr.NewHistoryServiceServer(historySvc)
	pb.RegisterHistoryServiceServer(grpcServer, historyGrpcSvc)

	// Schedule service
	scheduleRepo := repos.NewScheduleRepo(ds)
	scheduleSvc := services.NewScheduleService(scheduleRepo)
	scheduleGrpcSvc := grpcsvr.NewScheduleServer(scheduleSvc)
	pb.RegisterScheduleServiceServer(grpcServer, scheduleGrpcSvc)

	// Booking service
	bookingRepo := repos.NewBookingRepo(ds)
	bookingSvc := services.NewBookingService(bookingRepo)
	bookingGrpcSvc := grpcsvr.NewBookingServer(bookingSvc, scheduleSvc)
	pb.RegisterBookingServiceServer(grpcServer, bookingGrpcSvc)

	// Payment service
	paymentSvc := services.NewPaymentService(bookingRepo, nil, config.BaseURL)
	paymentGrpcSvc := grpcsvr.NewPaymentServer(paymentSvc)
	pb.RegisterPaymentServiceServer(grpcServer, paymentGrpcSvc)

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

	gwCtx, gwCancel := context.WithCancel(context.Background())
	defer gwCancel()
	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch strings.ToLower(key) {
			case "x-cart-session":
				return key, true
			default:
				return runtime.DefaultHeaderMatcher(key)
			}
		}),
	)
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

	err = gw.RegisterSMTPServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register SMTPService gateway")
	}

	err = gw.RegisterCatalogueServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register CatalogueService gateway")
	}

	err = gw.RegisterCartServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register CartService gateway")
	}

	err = gw.RegisterCustomerServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register CustomerService gateway")
	}

	err = gw.RegisterHistoryServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register HistoryService gateway")
	}

	err = gw.RegisterBookingServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register BookingService gateway")
	}

	err = gw.RegisterPaymentServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register PaymentService gateway")
	}

	err = gw.RegisterScheduleServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register ScheduleService gateway")
	}

	// ========================================
	// HTTP Server with Gateway + Chi
	// ========================================

	httpAddr := fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port)
	logStartupInfo(grpcServer, grpcAddr, httpAddr)

	server := thttp.NewServerWithGateway(config, healthSvc, authMiddleware, gwmux)
	err = server.Serve()

	if err != nil {
		return err
	}

	// Graceful shutdown of gRPC server
	grpcServer.GracefulStop()
	log.Info().Msg("gRPC server stopped gracefully")

	return nil
}

// logStartupInfo logs server addresses and all registered HTTP endpoints
// by reading the google.api.http annotations from the proto descriptors.
func logStartupInfo(grpcServer *grpc.Server, grpcAddr, httpAddr string) {
	log.Info().
		Str("http", httpAddr).
		Str("grpc", grpcAddr).
		Str("gateway_target", grpcAddr).
		Msg("server addresses")

	type endpoint struct {
		method  string
		path    string
		service string
		rpc     string
	}

	var endpoints []endpoint

	for serviceName, serviceInfo := range grpcServer.GetServiceInfo() {
		desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(serviceName))
		if err != nil {
			for _, m := range serviceInfo.Methods {
				log.Warn().
					Str("service", serviceName).
					Str("rpc", m.Name).
					Msg("no proto descriptor found, cannot resolve HTTP path")
			}
			continue
		}

		serviceDesc, ok := desc.(protoreflect.ServiceDescriptor)
		if !ok {
			continue
		}

		methods := serviceDesc.Methods()
		for i := 0; i < methods.Len(); i++ {
			md := methods.Get(i)
			opts := md.Options()
			if opts == nil {
				continue
			}

			httpRule, ok := proto.GetExtension(opts, annotations.E_Http).(*annotations.HttpRule)
			if !ok || httpRule == nil {
				continue
			}

			var method, path string
			switch p := httpRule.Pattern.(type) {
			case *annotations.HttpRule_Get:
				method, path = "GET", p.Get
			case *annotations.HttpRule_Post:
				method, path = "POST", p.Post
			case *annotations.HttpRule_Put:
				method, path = "PUT", p.Put
			case *annotations.HttpRule_Delete:
				method, path = "DELETE", p.Delete
			case *annotations.HttpRule_Patch:
				method, path = "PATCH", p.Patch
			}

			if path != "" {
				endpoints = append(endpoints, endpoint{
					method:  method,
					path:    path,
					service: serviceName,
					rpc:     string(md.Name()),
				})
			}
		}
	}

	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].path == endpoints[j].path {
			return endpoints[i].method < endpoints[j].method
		}
		return endpoints[i].path < endpoints[j].path
	})

	for _, ep := range endpoints {
		log.Info().
			Str("method", ep.method).
			Str("path", ep.path).
			Str("service", ep.service).
			Str("rpc", ep.rpc).
			Msg("endpoint")
	}

	log.Info().
		Str("method", "GET").Str("path", "/_internal/health").Msg("endpoint")
	log.Info().
		Str("method", "GET").Str("path", "/_internal/ready").Msg("endpoint")

	log.Info().Int("total", len(endpoints)+2).Msg("endpoints registered")
}
