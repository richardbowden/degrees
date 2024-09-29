package datastore

import (
	db "github.com/typewriterco/p402/internal/db_base"
)

type DataStorer interface {
	db.Storer
}

type dataStore struct {
	db.Storer
}

func NewDataStore(q db.Storer) DataStorer {
	u := dataStore{q}
	return u
}
