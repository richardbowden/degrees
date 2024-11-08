package main

import (
	"encoding/json"
	"net/http"

	"gitea.com/go-chi/binding"
	"github.com/go-chi/chi/v5"
	apihttp "github.com/typewriterco/p402/internal/api/http"
)

const (
	JSON_CONTENT_TYPE = "application/json; charset=utf-8"
)

func errorHandler(errs binding.Errors, rw http.ResponseWriter) {
	if len(errs) > 0 {
		rw.Header().Set("Content-Type", JSON_CONTENT_TYPE)
		if errs.Has(binding.ERR_DESERIALIZATION) {
			rw.WriteHeader(http.StatusBadRequest)
		} else if errs.Has(binding.ERR_CONTENT_TYPE) {
			rw.WriteHeader(http.StatusUnsupportedMediaType)
		} else {
			rw.WriteHeader(binding.STATUS_UNPROCESSABLE_ENTITY)
		}
		errOutput, _ := json.Marshal(errs)
		rw.Write(errOutput)
		return
	}
}

func (a *server) Endpoints() http.Handler {
	r := chi.NewRouter()
	// r.Get("/register", e.RegisterUserWithPassword)

	r.Post("/account/login", a.accountHandler.Login)

	r.Group(func(r chi.Router) {
		r.Post("/account/new", a.accountHandler.NewAccount)
	})

	r.Group(func(r chi.Router) {
		r.Use(apihttp.LogInMiddleware())
		r.Post("/account/logout", a.accountHandler.Logout)
	})
	r.Get("/debug", a.debugHandler.Debug)
	// r.Get("/user/sign-up/checkusername", a.CheckUsernameAvailability)

	return r
}
