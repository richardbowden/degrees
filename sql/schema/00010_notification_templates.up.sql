CREATE TABLE template (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    ref VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    scope_type TEXT NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),
    UNIQUE (ref, version)
);

SELECT add_updated_at_trigger('template');

INSERT INTO template (name, ref, content, scope_type, version, created_by, updated_by)
VALUES 
  ('Welcome Email', 'welcome-email', 'Dear {{.Name}}, welcome to our platform!', 'System', 1, NULL, NULL),
  ('Verify Email Address', 'verify-email-address', 'Click here to verify email address: {{.EmailVerifyURL}}', 'System', 1, NULL, NULL),
  ('Password Reset', 'password-reset', 'Click here to reset: {{.ResetLink}}', 'System', 1, NULL, NULL)
ON CONFLICT (ref, version) DO NOTHING;

CREATE TABLE notification_template (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    template_id BIGINT NOT NULL REFERENCES template(id)
);

INSERT INTO notification_template (name, template_id)
VALUES
    ('system-welcome-email', 1),
    ('system-verify-email-address', 2),
    ('system-password-reset', 2)
ON conflict (name) do nothing;
