package templater

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"sync"

	"github.com/typewriterco/p402/internal/dbpg"
)

type Template struct {
	ID        int64
	Name      string
	Ref       string
	Content   string
	Version   int
	ScopeType string
	CreatedBy string
	UpdatedBy string
}

type TemplateManager struct {
	store dbpg.Storer

	mu sync.RWMutex
	t  *template.Template
}

func NewTemplateManager(ctx context.Context, store dbpg.Store) (*TemplateManager, error) {

	allTemplats, err := store.ListSystemNotificationTemplates(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list system notificaiton templates %w", err)

	}

	t := template.New("root")

	for _, tpl := range allTemplats {
		_, err := t.New(tpl.SystemName).Parse(tpl.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", tpl.SystemName, err)
		}
	}

	return &TemplateManager{
		store: &store,
		t:     t,
	}, nil
}

func (t *TemplateManager) SaveTemplate(tplt Template) error {
	// t.store.template

	return nil
}

func (t *TemplateManager) GetTemplateByRef(ref string) (*Template, error) {
	return nil, nil
}

// func (t *TemplateManager) ListTemplates(ctx context.Context) error {
//
// 	allTemplates, err := t.store.ListSystemNotificdaionTemplates(ctx)
// }

func (ts *TemplateManager) RenderTemplate(ctx context.Context, templateName string, w io.Writer, data any) error {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.t.Lookup(templateName).Execute(w, data)
}
