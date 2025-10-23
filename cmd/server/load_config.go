package main

import (
	"fmt"
	"os"

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
)

func loadDBConfigFromCLI(ctx *cli.Context) config.DatabaseConfig {
	return config.DatabaseConfig{
		Host:     ctx.String(DBHostFlag),
		Port:     ctx.String(DBPortFlag),
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
		Database: loadDBConfigFromCLI(ctx),
		Auth:     loadAuthConfigFromCLI(ctx),
		SMTP: config.SMTPConfig{
			Host:     ctx.String(SMTPHostFlag),
			Port:     ctx.String(SMTPPortFlag),
			Username: ctx.String(SMTPUsernameFlag),
			Password: ctx.String(SMTPPasswordFlag),
		},
		DevOverrides: loadDevOverrides(),
	}
}

func loadDevOverrides() config.DevOverrideConfig {

	devOverridesEnabled, exists := os.LookupEnv("OP_DEV_OVERRIDE_ENABLE")

	if (exists && devOverridesEnabled == "0") || !exists {
		return config.DevOverrideConfig{}
	}

	fmt.Println("***********************************************************************************")
	fmt.Println("**********************                                        *********************")
	fmt.Println("**********************    DEV_OVERRIDES have been enabled     *********************")
	fmt.Println("**********************                                        *********************")
	fmt.Println("***********************************************************************************")
	fmt.Println("")

	orc := config.DevOverrideConfig{}

	skipUserConfirm, exists := os.LookupEnv("OP_DEV_SKIP_USER_CONFIRM_EMAIL")

	fmt.Printf("OP_DEV_SKIP_USER_CONFIG_EMAIL: %s\n", skipUserConfirm)
	if exists && skipUserConfirm == "1" {
		orc.SkipUserConfirm = true
	}

	fmt.Println("***********************************************************************************")
	fmt.Println("")
	return orc
}
