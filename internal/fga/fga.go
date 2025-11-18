package fga

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/storage/postgres"
	"github.com/openfga/openfga/pkg/storage/sqlcommon"
	fgafs "github.com/typewriterco/p402/fga"
	"go.uber.org/zap"

	"github.com/openfga/openfga/pkg/logger"
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
	ddd := logger.WithFormat("json")
	// ddd := &logger.ZapLogger{zapLogger}
	fgaServer, err := server.NewServerWithOpts(
		server.WithDatastore(datastore),
		server.WithLogger(ddd),
	)
	if err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}

	log.Println("OpenFGA server created (in-process, shared DB config)")
	return &FGA{
		fs:     fgafs.FGSModels,
		server: fgaServer,
	}, nil

}

func (f *FGA) ListFiles() {

	// files, err := iofs.New(f.fs, ".")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println(files)

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
