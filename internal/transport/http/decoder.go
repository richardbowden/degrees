package thttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DecodeJSONBody[T any](r *http.Request) (*T, error) {
	var dst T

	// Set a reasonable limit on request body size (e.g., 1MB)
	r.Body = http.MaxBytesReader(nil, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Optional: reject unknown fields

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return nil, fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil, errors.New("request body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			return nil, fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return nil, fmt.Errorf("request body contains unknown field %s", fieldName)

		case errors.Is(err, io.EOF):
			return nil, errors.New("request body must not be empty")

		case err.Error() == "http: request body too large":
			return nil, errors.New("request body must not be larger than 1MB")

		default:
			return nil, err
		}
	}

	// Check for extra data after JSON
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return nil, errors.New("request body must only contain a single JSON object")
	}

	return &dst, nil
}

// Alternative version that allows configurable options
type DecodeOptions struct {
	DisallowUnknownFields bool
	MaxBodySize           int64
}

func DecodeJSONBodyWithOptions[T any](r *http.Request, opts DecodeOptions) (*T, error) {
	var dst T

	if opts.MaxBodySize == 0 {
		opts.MaxBodySize = 1048576 // Default 1MB
	}

	r.Body = http.MaxBytesReader(nil, r.Body, opts.MaxBodySize)

	dec := json.NewDecoder(r.Body)
	if opts.DisallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(&dst)
	if err != nil {
		// Same error handling as above...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return nil, fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil, errors.New("request body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			return nil, fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return nil, fmt.Errorf("request body contains unknown field %s", fieldName)
		case errors.Is(err, io.EOF):
			return nil, errors.New("request body must not be empty")
		case err.Error() == "http: request body too large":
			return nil, fmt.Errorf("request body must not be larger than %d bytes", opts.MaxBodySize)
		default:
			return nil, err
		}
	}

	return &dst, nil
}
