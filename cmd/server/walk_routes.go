package main

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (s *server) walkRoutes() {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// route = strings.Replace(route, "/*/", "/", -1)

		l := log.With().Str("method", method).Str("path", route).Logger()

		if l.GetLevel() == zerolog.DebugLevel {
			for _, mw := range middlewares {
				name := runtime.FuncForPC(reflect.ValueOf(mw).Pointer()).Name()
				l.Info().Str("middleware", name).Msg("")
			}
		} else {
			l.Info().Int("total_middlewares", len(middlewares)).Msg("")
		}

		return nil
	}

	if err := chi.Walk(s.router, walkFunc); err != nil {
		log.Error().Msgf("Logging err: %s\n", err.Error())
	}
}
