package dbpg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TXStore struct {
	*Queries
	tx pgx.Tx
}

func (s *TXStore) Commit(ctx context.Context) error {
	ce := s.tx.Commit(ctx)
	return ce
}

func (s *TXStore) Rollback(ctx context.Context) error {
	te := s.tx.Rollback(ctx)

	if te != nil && !errors.Is(te, pgx.ErrTxClosed) {
		return te
	}

	return nil
}

type Storer interface {
	Querier
	GetTX(ctx context.Context) (*TXStore, error)
}

type SchemaMigration struct {
	Version int64 `json:"version"`
	Dirty   bool  `json:"dirty"`
}

type Store struct {
	*Queries
	dbpg *pgxpool.Pool
	SchemaMigration
}

func NewStore(db *pgxpool.Pool) *Store {
	//todo(rich): tidy up getting db version in newstore

	query := `SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1`

	var migration SchemaMigration
	ctx := context.Background()
	err := db.QueryRow(ctx, query).Scan(&migration.Version, &migration.Dirty)
	if err != nil {
		panic("cannot get db version, we should not get here")
	}
	//if err != nil {
	//	return nil, fmt.Errorf("failed to query schema_migrations: %w", err)
	//}
	return &Store{
		New(db),
		db,
		migration,
	}
}

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

func NewStoreCreateCon(conStr string) (*Store, error) {
	//TODO(rich): name needs to come from main at some point
	con, err := NewConnection(conStr, "p402 0.0.1-alpha")

	if err != nil {
		return nil, err
	}

	s := NewStore(con)
	return s, nil
}

func IsErrNoRows(err error) bool {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return true
		}
	}
	return false
}

func StringToPGString(str string) pgtype.Text {
	s := pgtype.Text{}
	if str == "" {
		return s
	}

	s.String = str
	s.Valid = true

	return s
}

func (s *Store) GetTX(ctx context.Context) (*TXStore, error) {
	tx, err := s.dbpg.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create a tx %w", err)
	}

	q := New(tx)

	return &TXStore{
		q,
		tx,
	}, nil
}

func (s *Store) CheckDB(ctx context.Context) error {
	err := s.dbpg.Ping(ctx)

	if err != nil {
		panic(err)
	}

	return nil
}
