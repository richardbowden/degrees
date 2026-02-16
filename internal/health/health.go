package health

import (
	"context"
	"fmt"
)

// Checker represents a health check for a specific subsystem
type Checker interface {
	Check(ctx context.Context) error
	Name() string
}

// Status represents the health status of a specific service
type Status struct {
	Service string `json:"service"`
	Healthy bool   `json:"healthy"`
	Error   string `json:"error,omitempty"`
}

// Service aggregates multiple health checkers
type Service struct {
	checkers map[string]Checker
}

// NewService creates a new health service with the given checkers
func NewService(checkers ...Checker) *Service {
	checkersMap := make(map[string]Checker)
	for _, checker := range checkers {
		checkersMap[checker.Name()] = checker
	}
	return &Service{checkers: checkersMap}
}

// CheckAll runs all registered health checks and returns their statuses
func (s *Service) CheckAll(ctx context.Context) ([]Status, error) {
	statuses := make([]Status, 0, len(s.checkers))
	hasFailure := false

	for name, checker := range s.checkers {
		status := Status{
			Service: name,
			Healthy: true,
		}

		if err := checker.Check(ctx); err != nil {
			status.Healthy = false
			status.Error = err.Error()
			hasFailure = true
		}

		statuses = append(statuses, status)
	}

	if hasFailure {
		return statuses, fmt.Errorf("one or more health checks failed")
	}

	return statuses, nil
}

// CheckDatabase runs only the database health check
// This can be used at call sites that need to verify DB health before critical operations
func (s *Service) CheckDatabase(ctx context.Context) error {
	checker, exists := s.checkers["database"]
	if !exists {
		return fmt.Errorf("database checker not registered")
	}
	return checker.Check(ctx)
}

// Check is a convenience method that returns true if a specific checker passes
func (s *Service) Check(ctx context.Context, name string) error {
	checker, exists := s.checkers[name]
	if !exists {
		return fmt.Errorf("checker %s not registered", name)
	}
	return checker.Check(ctx)
}
