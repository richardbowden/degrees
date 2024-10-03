package datastore

type DataStorer interface {
	Storer
}

type dataStore struct {
	Storer
}

func NewDataStore(q Storer) DataStorer {
	u := dataStore{q}
	return u
}
