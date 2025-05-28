package repos

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/services"
)

type Accounts struct {
	store dbpg.Storer
}

func NewAccountsRepo(store dbpg.Storer) *Accounts {
	return &Accounts{
		store: store}
}

func (a *Accounts) Create(ctx context.Context, params services.SignUpParams) (services.Account, error) {
	log := httplog.LogEntry(ctx)
	log.Info().Msg("from the account create repo layer")
	return services.Account{}, nil
}

func (a *Accounts) DoesAccountAlreadyExist(ctx context.Context, email string) (bool, error) {

	log := httplog.LogEntry(ctx)

	log.Info().Msg("this is from does account exist already")

	return true, nil

}
