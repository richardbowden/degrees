package fga

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/language/pkg/go/transformer"
	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/storage/postgres"
	"github.com/openfga/openfga/pkg/storage/sqlcommon"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"

	fgafs "github.com/typewriterco/p402/fga"
	"github.com/typewriterco/p402/internal/settings"
)

const (
	FGA_STORE_NAME = "p402"
)

type FGA struct {
	server *server.Server
	fs     embed.FS
	kv     settings.Settings

	storeID string
	authID  string
}

func New(ctx context.Context, pool *pgxpool.Pool, zerologger zerolog.Logger, kv settings.Settings) (*FGA, error) {

	datastore, err := postgres.NewWithDB(pool, pool, &sqlcommon.Config{
		MaxOpenConns:          25,
		MaxIdleConns:          5,
		ConnMaxIdleTime:       0,
		ConnMaxLifetime:       0,
		MaxTypesPerModelField: 100,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create fga datastore %w", err)
	}

	logger := NewZerologAdapter(zerologger)

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
		kv:     kv,
	}

	storeID, err := f.kv.Get(ctx, "fga_store_id")
	if err != nil {
		return &FGA{}, err
	}

	authID, err := f.kv.Get(ctx, "fga_store_id")
	if err != nil {
		return &FGA{}, err
	}

	if storeID == "" {
		storeID, authID, err = f.InitializeAuthModel(ctx)
		if err != nil {
			return &FGA{}, err
		}
	}

	log.Info().Msg("Creating OpenFGA store...")

	f.storeID = storeID
	f.authID = authID

	return f, nil

}

func (f *FGA) InitializeAuthModel(ctx context.Context) (string, string, error) {

	storeResp, err := f.server.CreateStore(ctx, &openfgav1.CreateStoreRequest{
		Name: FGA_STORE_NAME,
	})
	if err != nil {
		return "", "", fmt.Errorf("create store: %w", err)
	}
	storeID := storeResp.Id

	err = f.kv.Set(ctx, "fga_store_id", storeID)
	if err != nil {
		return "", "", fmt.Errorf("failed to stroe fga_store_id %w", err)
	}

	log.Printf("Store created: %s", storeID)

	modelBytes, err := embed.FS.ReadFile(f.fs, "models.fga")
	// fmt.Printf("$s\n", string(modelBytes))
	if err != nil {
		return "", "", fmt.Errorf("failed to read models.fga file from embdedd %w", err)
	}

	model, err := transformer.TransformDSLToProto(string(modelBytes))

	if err != nil {
		return "", "", fmt.Errorf("parse model: %w", err)
	}

	modelReq := &openfgav1.WriteAuthorizationModelRequest{
		SchemaVersion:   model.SchemaVersion,
		TypeDefinitions: model.TypeDefinitions,
		Conditions:      model.Conditions,
		StoreId:         storeID,
	}

	modelResp, err := f.server.WriteAuthorizationModel(ctx, modelReq)

	if err != nil {
		return "", "", fmt.Errorf("failed to write auth model %w", err)
	}

	err = f.kv.Set(ctx, "fga_mdoel_id", modelResp.AuthorizationModelId)

	if err != nil {
		return "", "", fmt.Errorf("failed to stroe fga_model_id %w", err)
	}

	return storeID, modelResp.AuthorizationModelId, nil

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
