package accesscontrol

import (
	"context"
	"crypto/sha256"
	"embed"
	"fmt"
	"io"
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

type AC struct {
	Server *server.Server
	fs     embed.FS
	kv     settings.Settings

	storeID string
	authID  string
}

func New(ctx context.Context, pool *pgxpool.Pool, zerologger zerolog.Logger, kv settings.Settings) (*AC, error) {

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
	f := &AC{
		fs:     fgafs.FGSModels,
		Server: fgaServer,
		kv:     kv,
	}

	log.Info().Msg("Creating OpenFGA store...")

	err = f.InitializeFGAModels(ctx)
	if err != nil {
		return &AC{}, err
	}

	return f, nil
}

func (f *AC) InitializeFGAModels(ctx context.Context) error {

	storeID, err := f.kv.Get(ctx, "fga_store_id")
	if err != nil {
		return err
	}

	if storeID == "" {
		storeResp, err := f.Server.CreateStore(ctx, &openfgav1.CreateStoreRequest{
			Name: FGA_STORE_NAME,
		})
		if err != nil {
			return fmt.Errorf("failed to create store: %w", err)
		}
		f.storeID = storeResp.Id
		f.kv.Set(ctx, "fga_store_id", f.storeID)
	}
	log.Printf("Store created: %s", storeID)

	r, _ := embed.FS.Open(f.fs, "models.fga")

	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		log.Error().Err(err).Msg("failed to hash fga model file")
		return fmt.Errorf("failed to hash fga model file: %w", err)
	}

	embdedModelFileHash := fmt.Sprintf("%x", h.Sum(nil))

	currentModelFileHash, err := f.kv.Get(ctx, "fga_model_file_hash")
	if err != nil {
		return fmt.Errorf("failed to fga model hash: %w", err)
	}

	fff := func() (string, error) {

		modelBytes, err := embed.FS.ReadFile(f.fs, "models.fga")
		if err != nil {
			return "", fmt.Errorf("failed to read models.fga file from embdedd %w", err)
		}

		authxModel, err := transformer.TransformDSLToProto(string(modelBytes))

		if err != nil {
			return "", fmt.Errorf("parse model: %w", err)
		}

		authxModelReq := &openfgav1.WriteAuthorizationModelRequest{
			SchemaVersion:   authxModel.SchemaVersion,
			TypeDefinitions: authxModel.TypeDefinitions,
			Conditions:      authxModel.Conditions,
			StoreId:         f.storeID,
		}

		authzModelResp, err := f.Server.WriteAuthorizationModel(ctx, authxModelReq)

		if err != nil {
			return "", fmt.Errorf("failed to write auth model %w", err)
		}
		return authzModelResp.GetAuthorizationModelId(), nil
	}

	writeModel := false
	if currentModelFileHash == "" {
		writeModel = true
	}

	//New version of the model
	if currentModelFileHash != embdedModelFileHash {
		writeModel = true
	}

	if writeModel {
		f.authID, err = fff()
		if err != nil {
			return fmt.Errorf("failed auth model stuff: %w", err)
		}

		err = f.kv.Set(ctx, "fga_model_file_hash", embdedModelFileHash)
		if err != nil {
			panic(err)
		}
		f.kv.Set(ctx, "fga_auth_id", f.authID)
		if err != nil {
			panic(err)
		}

	} else {
		authID, err := f.kv.Get(ctx, "fga_auth_id")
		if err != nil {
			return fmt.Errorf("failed to get auth model id from kv store: %w", err)
		}
		f.authID = authID
	}

	return nil
}

func (f *AC) ListFiles() {

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
