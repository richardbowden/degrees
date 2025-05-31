package repos

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v5/pgtype"
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

func (a *Accounts) Create(ctx context.Context, params services.PreparedSignupParams) (services.Account, error) {
	log := httplog.LogEntry(ctx)
	log.Info().Msg("from the account create repo layer")

	tx, err := a.store.GetTX(ctx)

	if err != nil {
		return services.Account{}, err
	}

	cap := dbpg.CreateAccountParams{
		FirstName:    params.FirstName,
		MiddleName:   pgtype.Text{String: params.MiddleName, Valid: true},
		Surname:      pgtype.Text{String: params.Surname, Valid: true},
		Email:        params.EMail,
		PasswordHash: params.HashedPassword,
		AccType:      dbpg.AccountTypeUser,
	}

	ac, err := tx.CreateAccount(ctx, cap)

	tx.Commit(ctx)
	_ = ac
	return services.Account{}, err
}

func (a *Accounts) DoesAccountAlreadyExist(ctx context.Context, email string) (bool, error) {
	return a.store.CheckAccountActiveEmailExists(ctx, dbpg.CheckAccountActiveEmailExistsParams{Email: email})
}
