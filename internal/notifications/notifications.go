package notification

import (
	"context"
	"strings"

	"github.com/go-chi/httplog"
	"github.com/typewriterco/p402/internal/riverqueue"
	"github.com/typewriterco/p402/internal/templater"
	"github.com/typewriterco/p402/internal/workers"
)

type Notifier struct {
	q         *riverqueue.RiverQueue
	tm        *templater.TemplateManager
	fromEmail string
}

func NewNotifier(q *riverqueue.RiverQueue, tm *templater.TemplateManager, fromEmail string) *Notifier {
	return &Notifier{
		q:         q,
		tm:        tm,
		fromEmail: fromEmail,
	}
}

// EmailData represents the common fields all email notifications need
type EmailData struct {
	To      []string
	Subject string
}

// SendEmail sends any email notification using the specified template
func (n *Notifier) SendEmail(ctx context.Context, templateType TemplateType, to []string, subject string, templateData any) error {
	log := httplog.LogEntry(ctx)

	var buf strings.Builder
	err := n.tm.RenderTemplate(ctx, templateType.String(), &buf, templateData)
	if err != nil {
		return err
	}

	log.Debug().
		Str("to", to[0]).
		Str("template", templateType.String()).
		Msg("sending email notification")

	emailJobArgs := workers.EmailArgs{
		To:      to,
		From:    n.fromEmail,
		Subject: subject,
		Content: buf.String(),
	}

	_, err = n.q.Client().Insert(ctx, emailJobArgs, nil)
	return err
}

// Template data structs for type safety at call sites
type VerifyEmailData struct {
	EmailVerifyURL string
}

type PasswordResetData struct {
	ResetLink string
	Email     string
}
