package notifications

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"sync"

	"github.com/typewriterco/p402/internal/dbpg"
)

type TemplateType string

const (
	TPL_WELCOME_EMAIL               TemplateType = "system-welcome-email"
	TPL_SYSTEM_VERIFY_EMAIL_ADDRESS TemplateType = "system-verify-email-address"
	TPL_SYSTEM_PASSWORD_RESET       TemplateType = "system-password-reset"
)

func (s TemplateType) String() string {
	return string(s)

}

type Template struct {
	ID      int          `json:"id"`
	Type    TemplateType `json:"type"`
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Content string       `json:"content"`
	Scope   string       `json:"scope"`
}

type TemplateService struct {
	store dbpg.Storer

	mu sync.RWMutex
	t  *template.Template
}

func NewTemplateService(ctx context.Context, store dbpg.Store) (*TemplateService, error) {

	allTemplats, err := store.ListSystemNotificationTemplates(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list system notificaiton templates %w", err)

	}

	t := template.New("root")

	for _, tpl := range allTemplats {
		t.New(tpl.SystemName).Parse(tpl.Content)
	}

	return &TemplateService{
		store: &store,
		t:     t,
	}, nil
}

func (ts *TemplateService) RenderTemplate(ctx context.Context, templateType TemplateType, w io.Writer, data any) error {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.t.Lookup(templateType.String()).Execute(w, data)
}
