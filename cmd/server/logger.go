package main

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func setBaseLogger(c *cli.Context) {

	humanizeLogs := c.Bool(HumanLogsFlag)
	logLevel := c.String("loglevel")

	logger := httplog.NewLogger(c.App.Name, httplog.Options{
		JSON:     !humanizeLogs,
		LogLevel: logLevel,
	}).With().Str("version", c.App.Version).Logger()

	log.Logger = logger
}

func (a *server)reportServerError(r *http.Request, err error) {
	method := r.Method
	url := r.URL.String()
	trace := string(debug.Stack())

	log.Error().Str("direction", "request").Str("method", method).Str("url", url).Str("trace", trace).Err(err).Msg("")
}


func (a *server)serverError(w http.ResponseWriter, r *http.Request, err error) {
 a.reportServerError(r, err)
	message := "The server encountered a problem and could not process your request"

	http.Error(w, message, http.StatusInternalServerError)
}
