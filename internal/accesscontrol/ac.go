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
	Server   *server.Server
	fs       embed.FS
	settings *settings.Service

	storeID string
	authID  string
}

func New(ctx context.Context, pool *pgxpool.Pool, zerologger zerolog.Logger, settingsService *settings.Service) (*AC, error) {

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
		fs:       fgafs.FGSModels,
		Server:   fgaServer,
		settings: settingsService,
	}

	log.Info().Msg("Creating OpenFGA store...")

	err = f.InitializeFGAModels(ctx)
	if err != nil {
		return &AC{}, err
	}

	return f, nil
}

func (f *AC) InitializeFGAModels(ctx context.Context) error {

	// Try to get existing store ID
	storeID, err := f.settings.GetString(ctx, "fga", "fga_store_id", settings.SystemScope())
	if err != nil && !settings.IsNotFound(err) {
		return fmt.Errorf("failed to get FGA store ID: %w", err)
	}

	if storeID == "" {
		// Create new store
		storeResp, err := f.Server.CreateStore(ctx, &openfgav1.CreateStoreRequest{
			Name: FGA_STORE_NAME,
		})
		if err != nil {
			return fmt.Errorf("failed to create store: %w", err)
		}
		f.storeID = storeResp.Id

		// Save store ID to settings
		if err = f.settings.SetSystem(ctx, "fga", "fga_store_id", f.storeID, nil, nil); err != nil {
			return fmt.Errorf("failed to save FGA store ID: %w", err)
		}
	} else {
		f.storeID = storeID
	}
	log.Info().Str("store_id", f.storeID).Msg("FGA store initialized")

	r, err := embed.FS.Open(f.fs, "models.fga")
	if err != nil {
		return fmt.Errorf("failed to open embedded fga model file: %w", err)
	}
	defer r.Close()

	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		log.Error().Err(err).Msg("failed to hash fga model file")
		return fmt.Errorf("failed to hash fga model file: %w", err)
	}

	embeddedModelFileHash := fmt.Sprintf("%x", h.Sum(nil))

	// Get current model hash
	currentModelFileHash, err := f.settings.GetString(ctx, "fga", "fga_model_file_hash", settings.SystemScope())
	if err != nil && !settings.IsNotFound(err) {
		return fmt.Errorf("failed to get FGA model hash: %w", err)
	}

	fff := func() (string, error) {

		modelBytes, err := embed.FS.ReadFile(f.fs, "models.fga")
		if err != nil {
			return "", fmt.Errorf("failed to read models.fga file from embedded %w", err)
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

	// Write model if hash is different or doesn't exist
	if currentModelFileHash != embeddedModelFileHash {
		f.authID, err = fff()
		if err != nil {
			return fmt.Errorf("failed to write authorization model: %w", err)
		}

		// Save model hash
		if err = f.settings.SetSystem(ctx, "fga", "fga_model_file_hash", embeddedModelFileHash, nil, nil); err != nil {
			return fmt.Errorf("failed to save FGA model hash: %w", err)
		}

		// Save auth model ID
		if err = f.settings.SetSystem(ctx, "fga", "fga_auth_id", f.authID, nil, nil); err != nil {
			return fmt.Errorf("failed to save FGA auth model ID: %w", err)
		}
	} else {
		// Use existing auth model ID
		authID, err := f.settings.GetString(ctx, "fga", "fga_auth_id", settings.SystemScope())
		if err != nil {
			return fmt.Errorf("failed to get FGA auth model ID: %w", err)
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

// GetStoreID returns the FGA store ID
func (f *AC) GetStoreID() string {
	return f.storeID
}

// GetAuthID returns the FGA authorization model ID
func (f *AC) GetAuthID() string {
	return f.authID
}

// Check verifies if a user has a specific relation to an object
func (f *AC) Check(ctx context.Context, user, relation, object string) (bool, error) {
	checkReq := &openfgav1.CheckRequest{
		StoreId:              f.storeID,
		AuthorizationModelId: f.authID,
		TupleKey: &openfgav1.CheckRequestTupleKey{
			User:     user,
			Relation: relation,
			Object:   object,
		},
	}

	checkResp, err := f.Server.Check(ctx, checkReq)
	if err != nil {
		return false, fmt.Errorf("fga check failed: %w", err)
	}

	return checkResp.GetAllowed(), nil
}

// WriteRelationship writes a single relationship tuple
func (f *AC) WriteRelationship(ctx context.Context, user, relation, object string) error {
	writeReq := &openfgav1.WriteRequest{
		StoreId:              f.storeID,
		AuthorizationModelId: f.authID,
		Writes: &openfgav1.WriteRequestWrites{
			TupleKeys: []*openfgav1.TupleKey{
				{
					User:     user,
					Relation: relation,
					Object:   object,
				},
			},
		},
	}

	_, err := f.Server.Write(ctx, writeReq)
	if err != nil {
		return fmt.Errorf("fga write failed: %w", err)
	}

	return nil
}

// DeleteRelationship removes a relationship tuple
func (f *AC) DeleteRelationship(ctx context.Context, user, relation, object string) error {
	deleteReq := &openfgav1.WriteRequest{
		StoreId:              f.storeID,
		AuthorizationModelId: f.authID,
		Deletes: &openfgav1.WriteRequestDeletes{
			TupleKeys: []*openfgav1.TupleKeyWithoutCondition{
				{
					User:     user,
					Relation: relation,
					Object:   object,
				},
			},
		},
	}

	_, err := f.Server.Write(ctx, deleteReq)
	if err != nil {
		return fmt.Errorf("fga delete failed: %w", err)
	}

	return nil
}
