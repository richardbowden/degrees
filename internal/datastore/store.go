package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typewriterco/p402/internal/dbpg"
)

type TXStore struct {
	*dbpg.Queries
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
	dbpg.Querier
	GetTX(ctx context.Context) (*TXStore, error)
}

type Store struct {
	*dbpg.Queries
	dbpg    *pgxpool.Pool
	SQLPool *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		dbpg.New(db),
		db,
		db,
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

func StringToSQLString(str string) sql.NullString {
	s := sql.NullString{}
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

	q := dbpg.New(tx)

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
