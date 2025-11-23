package settings

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/typewriterco/p402/internal/dbpg"
)

type Settings interface {
	Get(ctx context.Context, key string) (string, error)

	GetWithDefault(ctx context.Context, key, value string) (string, error)
	GetIntWithDefault(ctx context.Context, key string, value int) (int, error)

	Set(ctx context.Context, key, value string) error
	SetInt(ctx context.Context, key string, value int) error
}

type Backend interface {
	GetSetting(ctx context.Context, arg dbpg.GetSettingParams) (string, error)
	CreateSetting(ctx context.Context, arg dbpg.CreateSettingParams) error
}

type Store struct {
	store Backend
}

func New(db Backend) *Store {
	return &Store{store: db}
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	value, err := s.store.GetSetting(ctx, dbpg.GetSettingParams{Key: key})

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	return value, err
}

func (s *Store) GetWithDefault(ctx context.Context, key, value string) (string, error) {
	result, err := s.Get(ctx, key)

	if err != nil {
		return "", err
	}

	if result != "" {
		return result, nil
	}

	err = s.Set(ctx, key, value)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *Store) GetIntWithDefault(ctx context.Context, key string, value int) (int, error) {
	strValue := strconv.Itoa(value)

	result, err := s.GetWithDefault(ctx, key, strValue)
	if err != nil {
		return 0, err
	}

	intResult, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return intResult, nil
}

func (s *Store) Set(ctx context.Context, key, value string) error {
	params := dbpg.CreateSettingParams{
		Key:   key,
		Value: value,
	}
	return s.store.CreateSetting(ctx, params)
}

func (s *Store) SetInt(ctx context.Context, key string, value int) error {
	valueAsStr := strconv.Itoa(value)
	return s.Set(ctx, key, valueAsStr)
}

type TXSettingsStore struct {
	*Store
}

// func NewTXSettingStore(tx *dbpg.TXStore) TXSettingsStore {
// 	n := New(tx)
// 	i := TXSettingsStore{n}
// 	return i
// }
