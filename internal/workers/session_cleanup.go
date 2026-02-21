package workers

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/richardbowden/degrees/internal/dbpg"
)

const QueueMaintenance = "maintenance"

type SessionCleanupArgs struct{}

func (SessionCleanupArgs) Kind() string { return "session_cleanup" }

func (SessionCleanupArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: QueueMaintenance,
	}
}

type SessionCleanupWorker struct {
	river.WorkerDefaults[SessionCleanupArgs]
	store *dbpg.Store
}

func NewSessionCleanupWorker(store *dbpg.Store) *SessionCleanupWorker {
	return &SessionCleanupWorker{store: store}
}

func (w *SessionCleanupWorker) Work(ctx context.Context, job *river.Job[SessionCleanupArgs]) error {
	log.Info().Msg("starting session cleanup")

	err := w.store.DeleteExpiredSessions(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete expired sessions")
		return err
	}

	log.Info().Msg("session cleanup completed successfully")
	return nil
}
