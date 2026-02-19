package main

import (
	"github.com/typewriterco/p402/internal/config"
	"github.com/urfave/cli/v2"
)

const (
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
