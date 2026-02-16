-- Revert password reset template mapping
UPDATE notification_template
SET template_id = 2
WHERE name = 'system-password-reset';
