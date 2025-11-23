package config

import (
	"fmt"
)

type HTTPConfig struct {
	Port int
	Host string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type DatabaseConfig struct {
	Host                 string
	Port                 string
	User                 string
	Password             string
	DBName               string
	SSLMode              bool
	ConnectionRetryCount int
}

func (d *DatabaseConfig) ConnectionString() string {
	return d.ConnectionStringWithSchema("public")
}

func (d *DatabaseConfig) ConnectionStringWithSchema(schema string) string {
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

	c = fmt.Sprintf("%s&search_path=%s", c, schema)

	return c
}

type AuthConfig struct {
	GoogleClientID  string
	GoogleSecretKey string
	HostedDomain    string
	CookieLifetime  int
}

type DevOverrideConfig struct {
	SkipUserConfirm bool
}

type Config struct {
	Version          string
	HTTP             HTTPConfig
	Database         DatabaseConfig
	GoogleClientID   string
	GoogleSecKey     string
	HTTPPort         string
	HostedDomainName string
	CookieLifeTime   int
	Auth             AuthConfig
	DevOverrides     DevOverrideConfig
	SMTP             SMTPConfig
}
