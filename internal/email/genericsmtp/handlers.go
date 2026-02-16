//go:generate go tool valforge -file $GOFILE
package fastmail

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/httpjson"
	"github.com/typewriterco/p402/internal/problems"
	"github.com/typewriterco/p402/internal/valgen"
)

// ConfigureSMTPRequest represents the request body for configuring SMTP settings
type ConfigureSMTPRequest struct {
	SMTPAddress string `json:"smtp_address" validate:"required,minlen=3"`
	SMTPPort    int    `json:"smtp_port" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Identity    string `json:"identity"`
}

// ConfigureSMTPResponse represents the response after configuring SMTP
type ConfigureSMTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SMTPStatusResponse represents the current SMTP configuration status
type SMTPStatusResponse struct {
	Ready       bool   `json:"ready"`
	SMTPAddress string `json:"smtp_address,omitempty"`
	SMTPPort    int    `json:"smtp_port,omitempty"`
	Username    string `json:"username,omitempty"`
	Configured  bool   `json:"configured"`
}

// HTTPConfigureHandler handles POST requests to configure SMTP settings
func (c *Client) HTTPConfigureHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := httplog.LogEntry(ctx)

	// Decode request body
	req, err := httpjson.DecodeJSONBody[ConfigureSMTPRequest](r)
	if err != nil {
		log.Error().Err(err).Msg("problem parsing the json body")
		p := problems.New(problems.InvalidRequest, "problem parsing the json body", err)
		problems.WriteHTTPError(w, p)
		return
	}

	// Validate request
	err = req.Validate()
	var vErr *valgen.ValidationError
	if errors.As(err, &vErr) {
		p := problems.New(problems.InvalidRequest, "validation errors")
		p.AddDetails(vErr.Errors)
		problems.WriteHTTPError(w, p)
		return
	}

	// Update SMTP configuration
	err = c.SetConfig(ctx, req.SMTPAddress, req.SMTPPort, req.Username, req.Password, req.Identity)
	if err != nil {
		log.Error().Err(err).Msg("failed to configure SMTP")
		problems.WriteHTTPErrorWithErr(w, err)
		return
	}

	log.Info().
		Str("smtp_address", req.SMTPAddress).
		Int("smtp_port", req.SMTPPort).
		Str("username", req.Username).
		Msg("SMTP configuration updated successfully")

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ConfigureSMTPResponse{
		Success: true,
		Message: "SMTP configuration updated successfully",
	})
}

// HTTPStatusHandler handles GET requests to retrieve SMTP status
func (c *Client) HTTPStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := httplog.LogEntry(ctx)

	log.Info().Bool("ready", c.ready).Msg("SMTP status requested")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SMTPStatusResponse{
		Ready:       c.ready,
		SMTPAddress: c.config.SMTPAddress,
		SMTPPort:    c.config.SMTPPort,
		Username:    c.config.Username,
		Configured:  c.ready,
	})
}
