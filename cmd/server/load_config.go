package main

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/richardbowden/degrees/internal/config"
	"github.com/urfave/cli/v2"
)

// parseDatabaseURL parses a postgres:// URL into a DatabaseConfig.
// Used when DATABASE_URL is set (e.g. by fly postgres attach).
func parseDatabaseURL(rawURL string) (config.DatabaseConfig, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return config.DatabaseConfig{}, err
	}
	port, _ := strconv.Atoi(u.Port())
	if port == 0 {
		port = 5432
	}
	password, _ := u.User.Password()
	dbName := strings.TrimPrefix(u.Path, "/")
	sslMode := u.Query().Get("sslmode") == "verify-full"
	return config.DatabaseConfig{
		Host:     u.Hostname(),
		Port:     port,
		User:     u.User.Username(),
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}, nil
}

const (
	DatabaseURLFlag     = "database-url"
	DBHostFlag          = "db-host"
	DBPortFlag          = "db-port"
	DBUserFlag          = "db-user"
	DBNameFlag          = "db-name"
	DBPasswordFlag      = "db-pass"
	DBSSLModeFlag       = "db-sslmode"
	HTTPPortFlag        = "http-port"
	HTTPHostFlag        = "http-host"
	GRPCPortFlag        = "grpc-port"
	GRPCHostFlag        = "grpc-host"
	HumanLogsFlag       = "human-logs"
	LoggingLevelFlag    = "logging-level"
	GoogleClientIDFlag  = "google-client-id"  //TODO(rich): prob not needed here, will be moved into run time config
	GoogleClientSecFlag = "google-secret-key" //TODO(rich): prob not needed here, will be moved into run time config
	HostedDomainName    = "hosted-domain-name"
	CookieLifeTime      = "cookie-lifetime"
	SMTPHostFlag        = "smtp-host"
	SMTPPortFlag        = "smtp-port"
	SMTPUsernameFlag    = "smtp-username"
	SMTPPasswordFlag    = "smtp-password"
	FGAStoreIDFlag      = "fga-store-id"
	BaseURLFlag         = "base-url"
	DefaultFromEmailFlag = "default-from-email"
)

func loadDBConfigFromCLI(ctx *cli.Context) config.DatabaseConfig {
	// --database-url / DATABASE_URL takes precedence over individual DB flags
	if dbURL := ctx.String(DatabaseURLFlag); dbURL != "" {
		if cfg, err := parseDatabaseURL(dbURL); err == nil {
			return cfg
		}
	}
	return config.DatabaseConfig{
		Host:     ctx.String(DBHostFlag),
		Port:     ctx.Int(DBPortFlag),
		User:     ctx.String(DBUserFlag),
		DBName:   ctx.String(DBNameFlag),
		SSLMode:  ctx.Bool(DBSSLModeFlag),
		Password: ctx.String(DBPasswordFlag),
	}
}

func loadAuthConfigFromCLI(ctx *cli.Context) config.AuthConfig {
	return config.AuthConfig{
		GoogleClientID:  ctx.String(GoogleClientIDFlag),
		GoogleSecretKey: ctx.String(GoogleClientSecFlag),
		HostedDomain:    ctx.String(HostedDomainName),
		CookieLifetime:  ctx.Int(CookieLifeTime),
	}
}

func loadConfigFromCLI(ctx *cli.Context) *config.Config {
	return &config.Config{
		HTTP: config.HTTPConfig{
			Port: ctx.Int(HTTPPortFlag),
			Host: ctx.String(HTTPHostFlag),
		},
		GRPC: config.GRPCConfig{
			Port: ctx.Int(GRPCPortFlag),
			Host: ctx.String(GRPCHostFlag),
		},
		Database: loadDBConfigFromCLI(ctx),
		Auth:     loadAuthConfigFromCLI(ctx),
		SMTP: config.SMTPConfig{
			Host:     ctx.String(SMTPHostFlag),
			Port:     ctx.String(SMTPPortFlag),
			Username: ctx.String(SMTPUsernameFlag),
			Password: ctx.String(SMTPPasswordFlag),
		},
		BaseURL:          ctx.String(BaseURLFlag),
		DefaultFromEmail: ctx.String(DefaultFromEmailFlag),
	}
}
