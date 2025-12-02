package repos

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/services"
)

type Users struct {
	store dbpg.Storer
}

func NewUserRepo(store dbpg.Storer) *Users {
	return &Users{
		store: store}
}

func (a *Users) Create(ctx context.Context, params services.NewUser) (services.User, error) {
	log := httplog.LogEntry(ctx)
	log.Info().Msg("from the account create repo layer")

	tx, err := a.store.GetTX(ctx)

	if err != nil {
		return services.User{}, err
	}
	cap := dbpg.CreateUserParams{
		FirstName:    params.FirstName,
		MiddleName:   dbpg.StringToPGString(params.MiddleName),
		Surname:      dbpg.StringToPGString(params.Surname),
		Username:     params.Username,
		LoginEmail:   params.EMail,
		PasswordHash: params.HashedPassword,
	}

	ac, err := tx.CreateUser(ctx, cap)

	userEmailParams := dbpg.CreateUserEmailParams{
		Email:      params.EMail,
		IsVerified: false,
		UserID:     ac.ID,
	}

	_, err = tx.CreateUserEmail(ctx, userEmailParams)

	err = tx.Commit(ctx)

	return services.User{}, err
}

func (a *Users) DoesUserExist(ctx context.Context, email string, username string) (emailExists, usernameExists bool, err error) {
	userState, err := a.store.UserExists(ctx, dbpg.UserExistsParams{
		LoginEmail: email,
		Username:   username,
	})

	emailExists = userState.EmailExists
	usernameExists = userState.UsernameExists
	return
}
