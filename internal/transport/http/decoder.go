package thttp

import (
	"net/http"

	"github.com/richardbowden/degrees/internal/httpjson"
)

// DecodeJSONBody wraps the shared httpjson decoder for backward compatibility
func DecodeJSONBody[T any](r *http.Request) (*T, error) {
	return httpjson.DecodeJSONBody[T](r)
}

// DecodeOptions wraps the shared httpjson options
type DecodeOptions = httpjson.DecodeOptions

// DecodeJSONBodyWithOptions wraps the shared httpjson decoder with options
func DecodeJSONBodyWithOptions[T any](r *http.Request, opts DecodeOptions) (*T, error) {
	return httpjson.DecodeJSONBodyWithOptions[T](r, opts)
}
