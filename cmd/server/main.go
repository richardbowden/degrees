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
		Name:    "p402 server",
		Version: "v0.0.1+alpha",
		Flags: []cli.Flag{
			&cli.GenericFlag{Name: LoggingLevelFlag, Value: &EnumFlag{Enum: []string{"TRACE", "DEBUG", "ERROR", "INFO"}, Default: "INFO"}, Aliases: []string{"l"}, EnvVars: []string{"P402_LOG_LEVEL"}},
			&cli.BoolFlag{Name: HumanLogsFlag, Value: false, EnvVars: []string{"P402_HUMAN_LOGS"}},
			&cli.StringFlag{Name: DBHostFlag, Value: "127.0.0.1", EnvVars: []string{"P402_DB_HOST"}},
			&cli.IntFlag{Name: DBPortFlag, Value: 5432, EnvVars: []string{"P402_DB_PORT"}},
			&cli.StringFlag{Name: DBUserFlag, Value: "p402", EnvVars: []string{"P402_DB_USER"}},
			&cli.StringFlag{Name: DBPasswordFlag, EnvVars: []string{"P402_DB_PASS"}},
			&cli.StringFlag{Name: DBNameFlag, Value: "p402", EnvVars: []string{"P402_DB_NAME"}},
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
					&cli.StringFlag{Name: HostedDomainName, Value: "localhost", EnvVars: []string{"P402_HOSTED_DOMAIN_NAME"}},
					&cli.IntFlag{Name: CookieLifeTime, Value: 10, EnvVars: []string{"P402_COOKIE_LIFETIME"}},
					&cli.StringFlag{Name: SMTPHostFlag, Value: "localhost", EnvVars: []string{"P402_SMTP_HOST"}},
					&cli.StringFlag{Name: SMTPPortFlag, Value: "2500", EnvVars: []string{"P402_SMTP_PORT"}},
					&cli.StringFlag{Name: SMTPUsernameFlag, Value: "anything", EnvVars: []string{"P402_SMTP_USERNAME"}},
					&cli.StringFlag{Name: SMTPPasswordFlag, Value: "anypassword", EnvVars: []string{"P402_SMTP_PASSWORD"}},
				},
				Subcommands: []*cli.Command{
					{
						Name:   "run",
						Action: serverRun,
						Flags: []cli.Flag{
							&cli.IntFlag{Name: HTTPPortFlag, Value: 3030, Usage: "port number for http", EnvVars: []string{"P402_HTTP_PORT"}},
							&cli.StringFlag{Name: HTTPHostFlag, Value: "localhost", Usage: "host or ip address to listen on, set to ':' to bind to all ip available ip addresses", EnvVars: []string{"P402_HTTP_HOST"}},
						},
					},
				},
			}, {
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to start app")
	}
}
