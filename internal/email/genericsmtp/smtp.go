package fastmail

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/richardbowden/degrees/internal/settings"
)

type Client struct {
	mu             sync.RWMutex
	configChanged  bool
	auth           smtp.Auth
	reloadSettings bool
	settingsService *settings.Service
	config         config
	ready          bool
}

type config struct {
	SMTPAddress string
	SMTPPort    int
	smtpServer  string
	Username    string
	Password    string
	Identity    string `json:"identity,omitempty"`
}

// Status represents the current SMTP configuration status
type Status struct {
	Ready       bool
	SMTPAddress string
	SMTPPort    int
	Username    string
	Configured  bool
}

var ErrSMTPNotReady = errors.New("smtp not ready")

// IsReady returns true if the SMTP client is configured and ready to send emails
func (c *Client) IsReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ready
}

// GetStatus returns the current SMTP configuration status
func (c *Client) GetStatus() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Status{
		Ready:       c.ready,
		SMTPAddress: c.config.SMTPAddress,
		SMTPPort:    c.config.SMTPPort,
		Username:    c.config.Username,
		Configured:  c.ready,
	}
}

func (c *Client) loadConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx := context.Background()

	// Use new hierarchical settings service (system scope for SMTP config)
	conf, err := settings.GetTyped[config](
		ctx,
		c.settingsService,
		"smtp",
		"config",
		settings.SystemScope(),
	)

	if err != nil {
		c.ready = false
		log.Warn().Str("subsystem", "smtp").Msg("smtp settings are not configured. email notifications will queue until smtp is configured")
		return nil // Not an error - just not configured yet
	}

	c.config = conf
	c.config.smtpServer = fmt.Sprintf("%s:%d", c.config.SMTPAddress, c.config.SMTPPort)
	c.configChanged = true
	c.ready = true

	log.Info().
		Str("subsystem", "smtp").
		Str("server", c.config.SMTPAddress).
		Int("port", c.config.SMTPPort).
		Str("username", c.config.Username).
		Msg("smtp config loaded successfully - ready to send emails")

	return nil
}

func NewClient(settingsService *settings.Service) *Client {
	c := &Client{
		settingsService: settingsService,
		ready:          false,
	}

	// Load config from settings - if not configured, client remains in not-ready state
	// This is not an error condition - SMTP can be configured later via admin endpoint
	if err := c.loadConfig(); err != nil {
		log.Error().Err(err).Str("subsystem", "smtp").Msg("failed to load smtp config")
	}

	return c
}

func (c *Client) getSMTPAuth() smtp.Auth {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.auth == nil || c.configChanged {
		c.auth = smtp.PlainAuth(c.config.Identity, c.config.Username, c.config.Password, c.config.SMTPAddress)
		c.configChanged = false
	}
	return c.auth
}

// SetConfig persists SMTP configuration to the settings store and reloads the client
func (c *Client) SetConfig(ctx context.Context, smtpAddress string, smtpPort int, username, password, identity string) error {
	cfg := config{
		SMTPAddress: smtpAddress,
		SMTPPort:    smtpPort,
		Username:    username,
		Password:    password,
		Identity:    identity,
	}

	// Persist to settings store (system scope)
	err := c.settingsService.SetSystem(ctx, "smtp", "config", cfg, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to save smtp config: %w", err)
	}

	// Reload config to update client state
	return c.loadConfig()
}

func (c *Client) Send(sender string, rcpt []string, subject, body string) error {
	c.mu.RLock()
	ready := c.ready
	smtpServer := c.config.smtpServer
	c.mu.RUnlock()

	if !ready {
		return ErrSMTPNotReady
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("From: %s\r\n", sender))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", rcpt[0]))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(body)
	buf.WriteString("\r\n")

	return smtp.SendMail(
		smtpServer,
		c.getSMTPAuth(),
		sender,
		rcpt,
		[]byte(buf.String()),
	)
}
