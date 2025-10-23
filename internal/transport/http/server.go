package thttp

import (
	"context"
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
	"github.com/rs/zerolog/log"

	"github.com/typewriterco/p402/internal/config"
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
	config     *config.Config
	wg         sync.WaitGroup
	httpServer *http.Server
	startTime  time.Time

	middleware map[string]Middleware

	handlers Handlers
}

func NewServer(cfg *config.Config, handlers *Handlers) *Server {
	return &Server{
		config:    cfg,
		handlers:  *handlers,
		startTime: time.Now().UTC(),
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
	rlog := log.With().Str("subsystem", "roter-setup").Logger()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(httplog.RequestLogger(log.Logger))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
	r.Use(middleware.AllowContentType("application/json"))

	//TODO(rich): make private
	//r.Get("/health", s.healthCheck)
	//r.Get("/ready", s.readinessCheck)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Post("/verify", s.handlers.Users.VerifyEmail)
				r.Post("/login", s.handlers.Users.Login)
				r.Post("/reset-password", s.handlers.Users.ResetPassword)
			})

			r.Group(func(r chi.Router) {
				r.Post("/logout", s.handlers.Users.Logout)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Post("/", s.handlers.Users.NewUser)
			})

			r.Route("/{id}", func(r chi.Router) {
				r.Post("/enable", s.handlers.Users.ResetPassword)
				r.Post("/disable", s.handlers.Users.ResetPassword)
				r.Post("/reset-password", s.handlers.Users.ResetPassword)
			})
		})

		r.Route("/admin", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.handlers.Users.ListAllUsers)
			})
		})

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
	log.Info().
		Str("address", s.httpServer.Addr).
		Dur("startup_time", time.Since(s.startTime)).
		Msg("Server started")

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
	// Check database connection
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// TODO: Add actual database ping here
	// if err := s.db.Ping(ctx); err != nil {
	//     w.WriteHeader(http.StatusServiceUnavailable)
	//     return
	// }

	_ = ctx
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}
