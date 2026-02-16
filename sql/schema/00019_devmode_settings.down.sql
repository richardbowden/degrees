-- Remove dev mode settings
DELETE FROM settings
WHERE scope = 'system'
  AND subsystem = 'devmode'
  AND key IN ('enabled', 'skip_email_verification', 'disable_rate_limits', 'allow_insecure_auth');
