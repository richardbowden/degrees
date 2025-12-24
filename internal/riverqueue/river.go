package riverqueue

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/rs/zerolog"
)

const (
	RIVER_SCHEMA_NAME = "river"
)

type RiverQueue struct {
	Client *river.Client[pgx.Tx]
	pool   *pgxpool.Pool
}

// New creates a new River queue client
// This is not yet used in the application, but provides the foundation
// for future background job processing
func New(ctx context.Context, pool *pgxpool.Pool, logger zerolog.Logger) (*RiverQueue, error) {
	workers := river.NewWorkers()

	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Logger: NewZerologAdapter(logger),
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: workers,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	return &RiverQueue{
		Client: riverClient,
		pool:   pool,
	}, nil
}

// Start begins processing jobs from the queue
func (r *RiverQueue) Start(ctx context.Context) error {
	return r.Client.Start(ctx)
}

// Stop gracefully shuts down the River queue
func (r *RiverQueue) Stop(ctx context.Context) error {
	return r.Client.Stop(ctx)
}
