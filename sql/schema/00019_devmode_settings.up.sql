-- Insert default dev mode settings (all disabled by default)
-- Note: value column is JSONB, so booleans are JSON booleans (not strings)
INSERT INTO settings (scope, subsystem, key, value, description)
VALUES
    ('system', 'devmode', 'enabled', 'false', 'Enable development mode (disables security features for easier testing)'),
    ('system', 'devmode', 'skip_email_verification', 'false', 'Skip email verification in dev mode'),
    ('system', 'devmode', 'disable_rate_limits', 'false', 'Disable rate limiting in dev mode'),
    ('system', 'devmode', 'allow_insecure_auth', 'false', 'Allow insecure auth methods in dev mode');
