package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/richardbowden/degrees/internal/dbpg"
)

type Settings interface {
	Get(ctx context.Context, subsystem, key string) (string, error)

	GetWithDefault(ctx context.Context, subsystem, key, value string) (string, error)
	GetIntWithDefault(ctx context.Context, subsystem, key string, value int) (int, error)

	Set(ctx context.Context, subsystem, key, value string) error
	SetInt(ctx context.Context, subsystem, key string, value int) error
}

type Backend interface {
	GetSetting(ctx context.Context, arg dbpg.GetSettingParams) ([]byte, error)
	CreateSetting(ctx context.Context, arg dbpg.CreateSettingParams) error
}

type Store struct {
	store    Backend
	register map[string]any
}

func New(db Backend) *Store {
	return &Store{store: db}
}

func (s *Store) Get(ctx context.Context, subsystem, key string) (string, error) {
	if subsystem == "" {
		return "", fmt.Errorf("subsystem must not be empty")
	}
	value, err := s.store.GetSetting(ctx, dbpg.GetSettingParams{Subsystem: subsystem, Key: key})

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	return string(value), err
}

func (s *Store) GetWithDefault(ctx context.Context, subsystem, key, value string) (string, error) {
	result, err := s.Get(ctx, subsystem, key)

	if err != nil {
		return "", err
	}

	if result != "" {
		return result, nil
	}

	err = s.Set(ctx, subsystem, key, value)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *Store) GetIntWithDefault(ctx context.Context, subsystem, key string, value int) (int, error) {
	strValue := strconv.Itoa(value)

	result, err := s.GetWithDefault(ctx, subsystem, key, strValue)
	if err != nil {
		return 0, err
	}

	intResult, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return intResult, nil
}

func (s *Store) Set(ctx context.Context, subsystem, key, value string) error {
	if subsystem == "" {
		return fmt.Errorf("subsystem must not be empty")
	}
	params := dbpg.CreateSettingParams{
		Subsystem: subsystem,
		Key:       key,
		Value:     []byte(value),
	}
	return s.store.CreateSetting(ctx, params)
}

func SetData(s *Store, ctx context.Context, subsystem, key string, data any) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return s.Set(ctx, subsystem, key, string(j))
}

var ErrKeyNotSet = errors.New("key not set")

func GetData[T any](s *Store, ctx context.Context, subsystem, key string) (T, error) {
	var d T
	rawData, err := s.store.GetSetting(ctx, dbpg.GetSettingParams{
		Subsystem: subsystem,
		Key:       key,
	})

	if dbpg.IsErrNoRows(err) {
		log.Warn().Str("subsystem", "settings_getdata").Str("subsystem", subsystem).Str("key", key).Msg("is not set")
		return d, ErrKeyNotSet
	}

	if err != nil {
		return d, err
	}

	err = json.Unmarshal([]byte(rawData), &d)
	if err != nil {
		return d, err
	}

	return d, nil
}

func (s *Store) SetInt(ctx context.Context, subsystem, key string, value int) error {
	valueAsStr := strconv.Itoa(value)
	return s.Set(ctx, subsystem, key, valueAsStr)
}

type TXSettingsStore struct {
	*Store
}

// func NewTXSettingStore(tx *dbpg.TXStore) TXSettingsStore {
// 	n := New(tx)
// 	i := TXSettingsStore{n}
// 	return i
// }
