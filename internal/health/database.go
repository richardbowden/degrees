package health

import (
	"context"
)

// DBStore defines the interface for database health checking
// This avoids direct dependency on dbpg.Store
type DBStore interface {
	CheckDB(ctx context.Context) error
}

// DatabaseChecker implements health checking for the database
type DatabaseChecker struct {
	store DBStore
}

// NewDatabaseChecker creates a new database health checker
func NewDatabaseChecker(store DBStore) *DatabaseChecker {
	return &DatabaseChecker{store: store}
}

// Check verifies the database is accessible
func (d *DatabaseChecker) Check(ctx context.Context) error {
	return d.store.CheckDB(ctx)
}

// Name returns the identifier for this checker
func (d *DatabaseChecker) Name() string {
	return "database"
}
