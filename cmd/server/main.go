package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type EnumFlag struct {
	Enum     []string
	Default  string
	selected string
}

func (e *EnumFlag) Set(value string) error {
	v := strings.ToUpper(value)
	for _, enum := range e.Enum {
		if enum == v {
			e.selected = v
			return nil
		}
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e EnumFlag) String() string {
	if e.selected == "" {
		return e.Default
	}

	return e.selected
}

func main() {
	run(os.Args)
}

func run(args []string) {
	app := &cli.App{
		Name:    "degrees server",
		Version: getVersion(),
		Flags: []cli.Flag{
			&cli.GenericFlag{Name: LoggingLevelFlag, Value: &EnumFlag{Enum: []string{"TRACE", "DEBUG", "ERROR", "INFO"}, Default: "INFO"}, Aliases: []string{"l"}, EnvVars: []string{"DEGREES_LOG_LEVEL"}},
			&cli.BoolFlag{Name: HumanLogsFlag, Value: false, EnvVars: []string{"DEGREES_HUMAN_LOGS"}},
			&cli.StringFlag{Name: DatabaseURLFlag, Usage: "Full postgres connection URL; overrides individual DB flags", EnvVars: []string{"DATABASE_URL"}},
			&cli.StringFlag{Name: DBHostFlag, Value: "127.0.0.1", EnvVars: []string{"DEGREES_DB_HOST"}},
			&cli.IntFlag{Name: DBPortFlag, Value: 5432, EnvVars: []string{"DEGREES_DB_PORT"}},
			&cli.StringFlag{Name: DBUserFlag, Value: "degrees", EnvVars: []string{"DEGREES_DB_USER"}},
			&cli.StringFlag{Name: DBPasswordFlag, EnvVars: []string{"DEGREES_DB_PASS"}},
			&cli.StringFlag{Name: DBNameFlag, Value: "degrees", EnvVars: []string{"DEGREES_DB_NAME"}},
			&cli.BoolFlag{Name: DBSSLModeFlag, Value: false, EnvVars: []string{"DEGREES_DB_SSL_MODE"}},
		},
		Commands: []*cli.Command{
			{
				Name:   "build-info",
				Action: buildInfo,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "extended", Value: false,
					},
				},
			},
			{
				Name: "server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: HostedDomainName, Value: "localhost", EnvVars: []string{"DEGREES_HOSTED_DOMAIN_NAME"}},
					&cli.IntFlag{Name: CookieLifeTime, Value: 10, EnvVars: []string{"DEGREES_COOKIE_LIFETIME"}},
					&cli.StringFlag{Name: SMTPHostFlag, Value: "localhost", EnvVars: []string{"DEGREES_SMTP_HOST"}},
					&cli.StringFlag{Name: SMTPPortFlag, Value: "2500", EnvVars: []string{"DEGREES_SMTP_PORT"}},
					&cli.StringFlag{Name: SMTPUsernameFlag, Value: "anything", EnvVars: []string{"DEGREES_SMTP_USERNAME"}},
					&cli.StringFlag{Name: SMTPPasswordFlag, Value: "anypassword", EnvVars: []string{"DEGREES_SMTP_PASSWORD"}},
					&cli.StringFlag{Name: BaseURLFlag, Value: "http://localhost:8080", Usage: "Base URL for the application (used in emails, etc.)", EnvVars: []string{"DEGREES_BASE_URL"}},
					&cli.StringFlag{Name: DefaultFromEmailFlag, Value: "noreply@localhost", Usage: "Default from email address for notifications", EnvVars: []string{"DEGREES_DEFAULT_FROM_EMAIL"}},
				},
				Subcommands: []*cli.Command{
					{
						Name:   "run",
						Action: serverRun,
						Flags: []cli.Flag{
							&cli.IntFlag{Name: HTTPPortFlag, Value: 8080, Usage: "port number for http", EnvVars: []string{"DEGREES_HTTP_PORT"}},
							&cli.StringFlag{Name: HTTPHostFlag, Value: "localhost", Usage: "host or ip address to listen on, set to ':' to bind to all ip available ip addresses", EnvVars: []string{"DEGREES_HTTP_HOST"}},
							&cli.IntFlag{Name: GRPCPortFlag, Value: 9090, Usage: "port number for gRPC", EnvVars: []string{"DEGREES_GRPC_PORT"}},
							&cli.StringFlag{Name: GRPCHostFlag, Value: "localhost", Usage: "host or ip address for gRPC server", EnvVars: []string{"DEGREES_GRPC_HOST"}},
							&cli.StringFlag{Name: FGAStoreIDFlag, EnvVars: []string{"DEGREES_FGA_STORE_ID"}},
						},
					},
				},
			},
			{
				Name:  "seed",
				Usage: "Seed the database with initial data",
				Subcommands: []*cli.Command{
					{
						Name:   "run",
						Usage:  "Run all seed data (schedule, catalogue, test user)",
						Action: seedRun,
					},
				},
			},
			{
				Name: "db",
				Subcommands: []*cli.Command{
					{
						Name: "migration",
						Subcommands: []*cli.Command{
							{
								Name:   "up",
								Action: dbMigrate,
							},
							{
								Name:   "version",
								Action: dbCurrentVersion,
							},
						},
					},
				},
			},
		},
	}

	if err := app.Run(args); err != nil {
		log.Fatal().Err(err).Msg("failed to start app")
	}
}
