package fga

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/storage/postgres"
	"github.com/openfga/openfga/pkg/storage/sqlcommon"
	log "github.com/rs/zerolog/log"

	fgafs "github.com/typewriterco/p402/fga"
	"go.uber.org/zap"
)

type FGA struct {
	server *server.Server
	fs     embed.FS
}

func New(ctx context.Context, pool *pgxpool.Pool, zapLogger *zap.Logger) (*FGA, error) {

	datastore, err := postgres.NewWithDB(pool, pool, &sqlcommon.Config{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxIdleTime: 0,
		ConnMaxLifetime: 0,
	})

	// Create logging adapter
	logger := NewZerologAdapter(log.Logger)

	fgaServer, err := server.NewServerWithOpts(
		server.WithDatastore(datastore),
		server.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}

	logger.Info("OpenFGA server created (in-process, shared DB config")
	f := &FGA{
		fs:     fgafs.FGSModels,
		server: fgaServer,
	}

	_, _, err = f.InitializeAuthModel(ctx)

	if err != nil {
		return &FGA{}, err
	}

	return f, nil

}

func (f *FGA) InitializeAuthModel(ctx context.Context) (string, string, error) {
	// Create store
	log.Info().Msg("Creating OpenFGA store...")
	storeResp, err := f.server.CreateStore(ctx, &openfgav1.CreateStoreRequest{
		Name: "todo-app",
	})
	if err != nil {
		return "", "", fmt.Errorf("create store: %w", err)
	}
	storeID := storeResp.Id
	log.Printf("Store created: %s", storeID)

	return "", "", nil

}

func (f *FGA) ListFiles() {

	entries, err := fs.ReadDir(f.fs, ".")
	if err != nil {
		fmt.Println("Error reading embedded directory:", err)
		return
	}

	fmt.Println("Files in embedded directory 'myfiles':")
	for _, entry := range entries {
		if !entry.IsDir() { // Check if it's a file
			fmt.Println("-", entry.Name())
		}
	}
}
