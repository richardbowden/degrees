package main

import (
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func setBaseLogger(c *cli.Context) {

	humanizeLogs := c.Bool(HumanLogsFlag)
	logLevel := c.String(LoggingLevelFlag)

	logger := httplog.NewLogger(c.App.Name, httplog.Options{
		JSON:     !humanizeLogs,
		LogLevel: logLevel,
	}).With().Str("version", c.App.Version).Logger()

	log.Logger = logger
}
