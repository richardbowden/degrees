package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	DBHostFlag     = "db-host"
	DBPortFlag     = "db-port"
	DBUserFlag     = "db-user"
	DBNameFlag     = "db-name"
	DBPasswordFlag = "db-pass"
	DBSSLModeFlag  = "db-sslmode"

	HTTPPortFlag     = "http-port"
	HumanLogsFlag    = "human-logs"
	LoggingLevelFlag = "logging-level"

	GoogleClientIDFlag  = "google-client-id"
	GoogleClientSecFlag = "google-secret-key"

	HostedDomainName = "hosted-domain-name"
	CookieLifeTime   = "cookie-lifetime"

	SMTPHostFlag     = "smtp-host"
	SMTPPortFlag     = "smtp-port"
	SMTPUsernameFlag = "smtp-username"
	SMTPPasswordFlag = "smtp-password"
)

type config struct {
	db               DBConfig
	googleClientID   string
	googleSecKey     string
	httpPort         string
	hostedDomainName string
	cookieLifeTime   int

	devOverrideConfig DevOverrideConfig

	smtp SMTP
}

type SMTP struct {
	Host     string
	Port     string
	Username string
	Password string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  bool
}

func (d *DBConfig) ConnectionString() string {
	var c string

	creds := d.User

	if d.Password != "" {
		creds = fmt.Sprintf("%s:%s", creds, d.Password)
	}

	c = fmt.Sprintf("postgres://%s@%s:%s/%s", creds, d.Host, d.Port, d.DBName)

	sslMode := "disable"

	if d.SSLMode {
		sslMode = "verify-full"
	}

	c = fmt.Sprintf("%s?&sslmode=%s", c, sslMode)

	return c
}

func DBConfigFromCTX(c *cli.Context) DBConfig {
	dbConfig := DBConfig{
		Host:     c.String(DBHostFlag),
		Port:     c.String(DBPortFlag),
		User:     c.String(DBUserFlag),
		DBName:   c.String(DBNameFlag),
		SSLMode:  c.Bool(DBSSLModeFlag),
		Password: c.String(DBPasswordFlag),
	}

	if dbConfig.DBName == "" {
		fmt.Printf("DB Namm needs to be set %s\n", DBNameFlag)
		os.Exit(1)
	}

	return dbConfig
}

type DevOverrideConfig struct {
	SkipUserConfirm bool
}

func check_for_dev_overrides() DevOverrideConfig {

	dev_overrides_enabled, exists := os.LookupEnv("OP_DEV_OVERRIDE_ENABLE")

	if (exists && dev_overrides_enabled == "0") || !exists {
		return DevOverrideConfig{}
	}

	fmt.Println("***********************************************************************************")
	fmt.Println("**********************                                        *********************")
	fmt.Println("**********************    DEV_OVERRIDES have been enabled     *********************")
	fmt.Println("**********************                                        *********************")
	fmt.Println("***********************************************************************************")
	fmt.Println("")

	orc := DevOverrideConfig{}

	skipUserConfirm, exists := os.LookupEnv("OP_DEV_SKIP_USER_CONFIRM_EMAIL")

	fmt.Printf("OP_DEV_SKIP_USER_CONFIG_EMAIL: %s\n", skipUserConfirm)
	if exists && skipUserConfirm == "1" {
		orc.SkipUserConfirm = true
	}

	fmt.Println("***********************************************************************************")
	fmt.Println("")
	return orc
}

func GetSMTPConfig(ctx *cli.Context) SMTP {
	s := SMTP{
		Host:     ctx.String(SMTPHostFlag),
		Port:     ctx.String(SMTPPortFlag),
		Username: ctx.String(SMTPUsernameFlag),
		Password: ctx.String(SMTPPasswordFlag),
	}

	return s
}

func GetConfig(ctx *cli.Context) config {
	con := config{
		db:               DBConfigFromCTX(ctx),
		googleClientID:   ctx.String(GoogleClientIDFlag),
		googleSecKey:     ctx.String(GoogleClientSecFlag),
		httpPort:         ctx.String(HTTPPortFlag),
		hostedDomainName: ctx.String(HostedDomainName),
		cookieLifeTime:   ctx.Int(CookieLifeTime),
		smtp:             GetSMTPConfig(ctx),
	}

	con.devOverrideConfig = check_for_dev_overrides()

	return con
}
