-- Fix the password reset template mapping to point to correct template
UPDATE notification_template
SET template_id = 3
WHERE name = 'system-password-reset';
