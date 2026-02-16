package settings

import (
	"context"
)

// DevMode provides helper methods for checking development mode settings
type DevMode struct {
	service *Service
}

// NewDevMode creates a new DevMode helper
func NewDevMode(service *Service) *DevMode {
	return &DevMode{service: service}
}

// IsEnabled checks if development mode is enabled at system level
func (d *DevMode) IsEnabled(ctx context.Context) bool {
	enabled, err := d.service.GetBool(ctx, "devmode", "enabled", SystemScope())
	if err != nil {
		return false // Default to disabled if setting doesn't exist
	}
	return enabled
}

// SkipEmailVerification checks if email verification should be skipped in dev mode
func (d *DevMode) SkipEmailVerification(ctx context.Context) bool {
	if !d.IsEnabled(ctx) {
		return false
	}
	skip, err := d.service.GetBool(ctx, "devmode", "skip_email_verification", SystemScope())
	if err != nil {
		return false
	}
	return skip
}

// DisableRateLimits checks if rate limiting should be disabled in dev mode
func (d *DevMode) DisableRateLimits(ctx context.Context) bool {
	if !d.IsEnabled(ctx) {
		return false
	}
	disable, err := d.service.GetBool(ctx, "devmode", "disable_rate_limits", SystemScope())
	if err != nil {
		return false
	}
	return disable
}

// AllowInsecureAuth checks if insecure auth methods should be allowed in dev mode
func (d *DevMode) AllowInsecureAuth(ctx context.Context) bool {
	if !d.IsEnabled(ctx) {
		return false
	}
	allow, err := d.service.GetBool(ctx, "devmode", "allow_insecure_auth", SystemScope())
	if err != nil {
		return false
	}
	return allow
}
