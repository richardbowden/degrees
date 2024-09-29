package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

type Store struct {
	*Queries
	db      *pgxpool.Pool
	SQLPool *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		New(db),
		db,
		db,
	}
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
	tx, err := s.db.Begin(ctx)
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
	err := s.db.Ping(ctx)

	if err != nil {
		panic(err)
	}

	return nil
}
