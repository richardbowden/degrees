package repos

import "github.com/richardbowden/degrees/internal/dbpg"

type Templates struct {
	store dbpg.Storer
}

func NewTemplateRepo(store dbpg.Storer) *Templates {
	return &Templates{
		store: store,
	}
}
