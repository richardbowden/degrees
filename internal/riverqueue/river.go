package riverqueue

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/rs/zerolog"
)

type WorkerConfig struct {
	Name       string
	Queue      string
	MaxWorkers int
}

type WorkerStatus struct {
	Name       string `json:"name"`
	Queue      string `json:"queue"`
	MaxWorkers int    `json:"max_workers"`
	Paused     bool   `json:"paused"`
}

type RiverQueue struct {
	pool   *pgxpool.Pool
	logger zerolog.Logger

	mu      sync.RWMutex
	client  *river.Client[pgx.Tx]
	workers *river.Workers
	configs []WorkerConfig
	started bool
}

func New(pool *pgxpool.Pool, logger zerolog.Logger) *RiverQueue {
	return &RiverQueue{
		pool:    pool,
		logger:  logger,
		workers: river.NewWorkers(),
	}
}

// Register adds a worker. Must be called before Start.
// If queue is empty, defaults to worker name.
// If maxWorkers is 0, defaults to 10.
func Register[T river.JobArgs](rq *RiverQueue, cfg WorkerConfig, worker river.Worker[T]) error {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if rq.started {
		return fmt.Errorf("cannot register worker after Start")
	}

	if cfg.Name == "" {
		return fmt.Errorf("worker name is required")
	}

	if cfg.Queue == "" {
		cfg.Queue = cfg.Name
	}

	if cfg.MaxWorkers == 0 {
		cfg.MaxWorkers = 10
	}

	river.AddWorker(rq.workers, worker)
	rq.configs = append(rq.configs, cfg)

	rq.logger.Debug().
		Str("worker", cfg.Name).
		Str("queue", cfg.Queue).
		Int("max_workers", cfg.MaxWorkers).
		Msg("registered worker")

	return nil
}

func (rq *RiverQueue) Start(ctx context.Context) error {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if rq.started {
		return fmt.Errorf("already started")
	}

	if len(rq.configs) == 0 {
		return fmt.Errorf("no workers registered")
	}

	queues := make(map[string]river.QueueConfig, len(rq.configs))
	for _, cfg := range rq.configs {
		queues[cfg.Queue] = river.QueueConfig{MaxWorkers: cfg.MaxWorkers}
	}

	client, err := river.NewClient(riverpgxv5.New(rq.pool), &river.Config{
		Logger:  NewZerologAdapter(rq.logger),
		Queues:  queues,
		Workers: rq.workers,
	})
	if err != nil {
		return fmt.Errorf("failed to create river client: %w", err)
	}

	rq.client = client
	rq.started = true

	return rq.client.Start(ctx)
}

func (rq *RiverQueue) Stop(ctx context.Context) error {
	rq.mu.RLock()
	client := rq.client
	rq.mu.RUnlock()

	if client == nil {
		return nil
	}
	return client.Stop(ctx)
}

// Client returns the underlying River client for inserting jobs.
// Returns nil if not started.
func (rq *RiverQueue) Client() *river.Client[pgx.Tx] {
	rq.mu.RLock()
	defer rq.mu.RUnlock()
	return rq.client
}

func (rq *RiverQueue) findQueue(name string) (string, error) {
	for _, cfg := range rq.configs {
		if cfg.Name == name {
			return cfg.Queue, nil
		}
	}
	return "", fmt.Errorf("worker %q not found", name)
}

func (rq *RiverQueue) Pause(ctx context.Context, name string) error {
	rq.mu.RLock()
	client := rq.client
	rq.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("not started")
	}

	queue, err := rq.findQueue(name)
	if err != nil {
		return err
	}

	return client.QueuePause(ctx, queue, nil)
}

func (rq *RiverQueue) Resume(ctx context.Context, name string) error {
	rq.mu.RLock()
	client := rq.client
	rq.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("not started")
	}

	queue, err := rq.findQueue(name)
	if err != nil {
		return err
	}

	return client.QueueResume(ctx, queue, nil)
}

func (rq *RiverQueue) List(ctx context.Context) ([]WorkerStatus, error) {
	rq.mu.RLock()
	client := rq.client
	configs := rq.configs
	rq.mu.RUnlock()

	if client == nil {
		return nil, fmt.Errorf("not started")
	}

	result := make([]WorkerStatus, 0, len(configs))
	for _, cfg := range configs {
		q, err := client.QueueGet(ctx, cfg.Queue)
		if err != nil {
			return nil, fmt.Errorf("failed to get queue %q: %w", cfg.Queue, err)
		}

		result = append(result, WorkerStatus{
			Name:       cfg.Name,
			Queue:      cfg.Queue,
			MaxWorkers: cfg.MaxWorkers,
			Paused:     q.PausedAt != nil,
		})
	}

	return result, nil
}
