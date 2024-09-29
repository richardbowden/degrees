package services

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/datastore"
)

type AccountService struct {
	ds datastore.DataStorer
}

func (ac *AccountService) NewAccount(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called new account... this does not work yet")

	return nil
}

func (ac *AccountService) Login(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called login.... this does not work yet")

	return nil
}

func (ac *AccountService) Logout(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called logout.. this does not work yet")

	return nil

}

func NewAccountService(ds datastore.DataStorer) (*AccountService, error) {

	ac := &AccountService{ds: ds}

	return ac, nil
}
