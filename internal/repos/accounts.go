package repos

import (
	"github.com/typewriterco/p402/internal/dbpg"
)

type Accounts struct {
	store dbpg.Storer
}

func NewAccountsRepo(store dbpg.Storer) *Accounts {
	return &Accounts{
		store: store}
}
