package thttp

import (
	"errors"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/problems"
	"github.com/typewriterco/p402/internal/services"
	"github.com/typewriterco/p402/internal/valgen"
)

type SignUp struct {
	svc *services.SignUp
}

func NewSignUp(signUpSvc services.SignUp) *SignUp {
	return &SignUp{svc: &signUpSvc}
}

func (s *SignUp) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := httplog.LogEntry(ctx)

	user, err := DecodeJSONBody[services.NewUserRequest](r)

	if err != nil {
		log.Error().Err(err).Msg("problem parsing the json body")
		p := problems.New(problems.InvalidRequest, "problem parsing the json body", err)
		problems.WriteHTTPError(w, p)
		return
	}

	err = user.Validate()

	var vErr *valgen.ValidationError
	if errors.As(err, &vErr) {
		p := problems.New(problems.InvalidRequest, "validation errors")
		p.AddDetails(vErr.Errors)
		problems.WriteHTTPError(w, p)
		return
	}

	err = s.svc.Register(ctx, user)

	if err != nil {
		log.Error().Err(err).Msg("")
		problems.WriteHTTPErrorWithErr(w, err)
		return
	}
}
