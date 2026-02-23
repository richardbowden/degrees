package repos

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/services"
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
	defer tx.Rollback(ctx)
	cap := dbpg.CreateUserParams{
		FirstName:    params.FirstName,
		MiddleName:   dbpg.StringToPGString(params.MiddleName),
		Surname:      dbpg.StringToPGString(params.Surname),
		Username:     params.Username,
		LoginEmail:   params.EMail,
		PasswordHash: params.HashedPassword,
		SignUpStage:  string(params.State),
	}

	createdUser, err := tx.CreateUser(ctx, cap)

	if err != nil {
		return services.User{}, err
	}

	userEmailParams := dbpg.CreateUserEmailParams{
		Email:      params.EMail,
		IsVerified: false,
		UserID:     createdUser.ID,
	}

	createdEmail, err := tx.CreateUserEmail(ctx, userEmailParams)

	if err != nil {
		return services.User{}, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return services.User{}, err
	}

	u := services.User{
		ID:         createdUser.ID,
		FirstName:  createdUser.FirstName,
		MiddleName: createdUser.MiddleName.String,
		Surname:    createdUser.Surname.String,
		EMail:      createdEmail.Email,
		// SignUpStage: createdUser.SignUpStage,
		Enabled:   createdUser.Enabled,
		Sysop:     createdUser.Sysop,
		CreatedOn: createdUser.CreatedOn.Time,
		UpdatedAt: createdUser.UpdatedAt.Time,
	}

	return u, nil
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

func (a *Users) UpdateSysop(ctx context.Context, userID int64, sysop bool) error {
	_, err := a.store.UpdateUserSysop(ctx, dbpg.UpdateUserSysopParams{
		ID:    userID,
		Sysop: sysop,
	})
	return err
}

func (a *Users) IsFirstUser(ctx context.Context) (bool, error) {
	return a.store.IsFirstUser(ctx)
}

func (a *Users) UpdateEnabled(ctx context.Context, userID int64, enabled bool) (services.User, error) {
	updatedUser, err := a.store.UpdateUserEnabled(ctx, dbpg.UpdateUserEnabledParams{
		ID:      userID,
		Enabled: enabled,
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.User{}, services.ErrNoRecord
		}
		return services.User{}, err
	}

	return services.User{
		ID:         updatedUser.ID,
		FirstName:  updatedUser.FirstName,
		MiddleName: updatedUser.MiddleName.String,
		Surname:    updatedUser.Surname.String,
		EMail:      updatedUser.LoginEmail,
		Enabled:    updatedUser.Enabled,
		Sysop:      updatedUser.Sysop,
		CreatedOn:  updatedUser.CreatedOn.Time,
		UpdatedAt:  updatedUser.UpdatedAt.Time,
	}, nil
}

func (a *Users) GetUserByID(ctx context.Context, userID int64) (services.User, error) {
	dbUser, err := a.store.GetUserById(ctx, dbpg.GetUserByIdParams{ID: userID})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.User{}, services.ErrNoRecord
		}
		return services.User{}, err
	}

	return services.User{
		ID:          dbUser.ID,
		FirstName:   dbUser.FirstName,
		MiddleName:  dbUser.MiddleName.String,
		Surname:     dbUser.Surname.String,
		EMail:       dbUser.LoginEmail,
		SignUpStage: dbUser.SignUpStage,
		Enabled:     dbUser.Enabled,
		Sysop:       dbUser.Sysop,
		CreatedOn:   dbUser.CreatedOn.Time,
		UpdatedAt:   dbUser.UpdatedAt.Time,
	}, nil
}

func (a *Users) UpdateUser(ctx context.Context, userID int64, firstName string, middleName string, surname string) (services.User, error) {
	dbUser, err := a.store.UpdateUser(ctx, dbpg.UpdateUserParams{
		ID:         userID,
		FirstName:  firstName,
		MiddleName: dbpg.StringToPGString(middleName),
		Surname:    dbpg.StringToPGString(surname),
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return services.User{}, services.ErrNoRecord
		}
		return services.User{}, err
	}

	return services.User{
		ID:          dbUser.ID,
		FirstName:   dbUser.FirstName,
		MiddleName:  dbUser.MiddleName.String,
		Surname:     dbUser.Surname.String,
		EMail:       dbUser.LoginEmail,
		SignUpStage: dbUser.SignUpStage,
		Enabled:     dbUser.Enabled,
		Sysop:       dbUser.Sysop,
		CreatedOn:   dbUser.CreatedOn.Time,
		UpdatedAt:   dbUser.UpdatedAt.Time,
	}, nil
}

func (a *Users) ListAllUsers(ctx context.Context) ([]services.User, error) {
	dbUsers, err := a.store.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]services.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = services.User{
			ID:         dbUser.ID,
			FirstName:  dbUser.FirstName,
			MiddleName: dbUser.MiddleName.String,
			Surname:    dbUser.Surname.String,
			EMail:      dbUser.LoginEmail,
			Enabled:    dbUser.Enabled,
			Sysop:      dbUser.Sysop,
			CreatedOn:  dbUser.CreatedOn.Time,
			UpdatedAt:  dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}
