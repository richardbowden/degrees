package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(conStr string, conName string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	conConfig, err := pgxpool.ParseConfig(conStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse db config %w", err)
	}

	conConfig.MaxConnIdleTime = time.Minute

	conConfig.ConnConfig.RuntimeParams["application_name"] = conName

	con, err := pgxpool.NewWithConfig(ctx, conConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool to database %w", err)
	}

	err = con.Ping(ctx)

	if err != nil {
		return nil, err
	}

	return con, nil
}
