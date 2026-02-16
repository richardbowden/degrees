package workers

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
)

type EMailer interface {
	Send(from string, rcpt []string, subject, body string) error
	IsReady() bool
}

const QueueEmail = "email"

type EmailArgs struct {
	To      []string `json:"to"`
	From    string   `json:"from"`
	Subject string   `json:"subject"`
	Content string   `json:"content"`

	// Optional callback info - what to notify on success
	CallbackType string `json:"callback_type,omitempty"` // e.g. "signup", "password_reset"
	CallbackID   string `json:"callback_id,omitempty"`   // e.g. signup ID, user ID
}

func (EmailArgs) Kind() string { return "email" }

func (EmailArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: QueueEmail}
}

// EmailCallback handles post-send actions
type EmailCallback interface {
	OnEmailSent(ctx context.Context, callbackType, callbackID string) error
}

type EmailWorker struct {
	river.WorkerDefaults[EmailArgs]
	mailer   EMailer
	callback EmailCallback // optional
}

// func NewEmailWorker(mailer Mailer, callback EmailCallback) *EmailWorker {
func NewEmailWorker(callback EmailCallback, mailer EMailer) *EmailWorker {
	return &EmailWorker{
		mailer:   mailer,
		callback: callback,
	}
}

func (w *EmailWorker) Work(ctx context.Context, job *river.Job[EmailArgs]) error {
	// Check if mailer is ready before attempting to process job
	if !w.mailer.IsReady() {
		// Return error so River will retry the job later
		// Jobs will remain in queue until SMTP is configured
		return fmt.Errorf("smtp not configured - job will retry once smtp is ready")
	}

	if len(job.Args.To) == 0 {
		return fmt.Errorf("no recipients")
	}

	err := w.mailer.Send(job.Args.From, job.Args.To, job.Args.Subject, job.Args.Content)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", job.Args.To[0], err)
	}

	// Fire callback if specified
	if w.callback != nil && job.Args.CallbackType != "" {
		if err := w.callback.OnEmailSent(ctx, job.Args.CallbackType, job.Args.CallbackID); err != nil {
			// Log but don't fail - email was sent successfully
			// Could also queue a separate job to retry the callback
			return fmt.Errorf("email sent but callback failed: %w", err)
		}
	}

	return nil
}
