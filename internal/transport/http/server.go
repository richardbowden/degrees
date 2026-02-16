package thttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"

	"github.com/typewriterco/p402/internal/config"
	"github.com/typewriterco/p402/internal/health"
)

var (
	IdleTimeout    = 1 * time.Minute
	ReadTimeout    = 5 * time.Second
	WriteTimeout   = 10 * time.Second
	ShutdownPeriod = 30 * time.Second

	AppName = "p402-backend"
)

type Middleware http.Handler

type Server struct {
	config         *config.Config
	healthSvc      *health.Service
	gatewayMux     *runtime.ServeMux
	wg             sync.WaitGroup
	httpServer     *http.Server
	startTime      time.Time
	authMiddleware *AuthMiddleware

	middleware map[string]Middleware

	handlers Handlers
}

func NewServer(cfg *config.Config, healthSvc *health.Service, handlers *Handlers, authMiddleware *AuthMiddleware) *Server {
	return &Server{
		config:         cfg,
		healthSvc:      healthSvc,
		handlers:       *handlers,
		authMiddleware: authMiddleware,
		startTime:      time.Now().UTC(),
	}
}

// NewServerWithGateway creates a server with gRPC-Gateway integration
func NewServerWithGateway(cfg *config.Config, healthSvc *health.Service, handlers *Handlers, authMiddleware *AuthMiddleware, gatewayMux *runtime.ServeMux) *Server {
	return &Server{
		config:         cfg,
		healthSvc:      healthSvc,
		gatewayMux:     gatewayMux,
		handlers:       *handlers,
		authMiddleware: authMiddleware,
		startTime:      time.Now().UTC(),
	}
}

func (s *Server) RegisterMIddleware(name string, middleware http.Handler) {
	s.middleware[name] = middleware
}

func (s *Server) Serve() error {
	router := s.setupRoutes()

	addr := fmt.Sprintf("%s:%d", s.config.HTTP.Host, s.config.HTTP.Port)

	s.httpServer = &http.Server{Addr: addr, Handler: router, IdleTimeout: IdleTimeout, ReadTimeout: ReadTimeout, WriteTimeout: WriteTimeout}

	return s.serveWithGracefulShutdown()
}

func (s *Server) setupRoutes() chi.Router {
	r := chi.NewMux()
	rlog := log.With().Str("subsystem", "router-setup").Logger()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(httplog.RequestLogger(log.Logger))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
	r.Use(middleware.AllowContentType("application/json"))

	// Internal endpoints for health checks
	r.Route("/_internal", func(r chi.Router) {
		r.Get("/health", s.healthCheck)
		r.Get("/ready", s.readinessCheck)
	})

	// Mount gRPC-Gateway
	// All API endpoints (/api/v1/*) are handled by gRPC-Gateway (auto-generated from proto)
	rlog.Info().Msg("mounting gRPC-Gateway at /api/v1")
	r.Mount("/", s.gatewayMux)

	// Admin routes for non-gRPC functionality
	r.Route("/admin", func(r chi.Router) {
		// Apply authentication and authorization middleware
		if s.authMiddleware != nil {
			r.Use(s.authMiddleware.RequireAuth)  // First authenticate
			r.Use(s.authMiddleware.RequireSysop) // Then check sysop role
		}

		// SMTP configuration endpoints (not in proto - direct HTTP only)
		if s.handlers.SMTP != nil {
			r.Route("/smtp", func(r chi.Router) {
				r.Post("/configure", s.handlers.SMTP.HTTPConfigureHandler)
				r.Get("/status", s.handlers.SMTP.HTTPStatusHandler)
			})
		}
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"xxxxendpoint not found"}`))
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		handlerName := getHandlerName(handler)
		middlewareCount := len(middlewares)
		rlog.Info().Str("method", method).Str("route", route).Str("handlerName", handlerName).Int("middleware(s)", middlewareCount).Msg("")
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		rlog.Error().Err(err).Msg("Error walking routes")
	}

	return r
}

func (s *Server) serveWithGracefulShutdown() error {
	shutdownChan := make(chan error, 1)

	// Handle shutdown signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info().Msg("Shutdown signal received")

		ctx, cancel := context.WithTimeout(context.Background(), ShutdownPeriod)
		defer cancel()

		shutdownChan <- s.httpServer.Shutdown(ctx)
	}()

	// Start server
	log.Info().Str("address", s.httpServer.Addr).Dur("startup_time", time.Since(s.startTime)).Msg("Server started")

	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}

	// Wait for shutdown
	err = <-shutdownChan
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	// Wait for all goroutines
	s.wg.Wait()

	log.Info().
		Dur("uptime", time.Since(s.startTime)).
		Msg("Server stopped gracefully")

	return nil
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func (s *Server) readinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	statuses, err := s.healthSvc.CheckAll(ctx)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Warn().Err(err).Msg("readiness check failed")
		w.WriteHeader(http.StatusServiceUnavailable)
		response := map[string]interface{}{
			"status":   "not ready",
			"services": statuses,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":   "ready",
		"services": statuses,
	}
	json.NewEncoder(w).Encode(response)
}
