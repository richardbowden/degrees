//go:generate go tool valforge -file $GOFILE
package services

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/richardbowden/passwordHash"
	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/problems"
)

type AuthN struct {
	db dbpg.Storer
}

func NewAuthN(db dbpg.Storer) *AuthN {
	return &AuthN{
		db: db,
	}
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

func (a *AuthN) Login(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called login.... this does not work yet")

	return nil
}

func (a *AuthN) Logout(ctx context.Context) error {

	l := httplog.LogEntry(ctx)

	l.Warn().Msg("called logout.. this does not work yet")

	return nil
}

func (a *AuthN) HashPassword(pwd1, pwd2 string) (string, error) {
	hashedPassword, err := passwordHash.HashWithDefaults(pwd1, pwd2)
	return hashedPassword, err
}

func (a *AuthN) SetOrUpdatePassword(ctx context.Context, userID int, pwd1, pwd2 string) error {
	hashedPassword, err := a.HashPassword(pwd1, pwd2)
	if err != nil {
		return err
	}
	_ = hashedPassword
	return problems.New(problems.Other, "setorupdate not done yet")

}
