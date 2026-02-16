//go:build tools
// +build tools

// Package tools tracks development tool dependencies.
// These tools will be included in go.mod but not in builds.
package tools

import (
	// Protocol buffer tooling
	_ "github.com/bufbuild/buf/cmd/buf"

	// gRPC testing tools (optional, comment out if not needed)
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
	_ "github.com/fullstorydev/grpcui/cmd/grpcui"

	// Existing tools
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)

// Usage:
//   go mod tidy           # Add tool dependencies to go.mod
//   go install <tool>     # Install a specific tool
//   go run <tool> [args]  # Run tool without installing
