-- Enhanced settings table for multi-tenant SaaS config
-- Supports hierarchical config: system → organization → project → user

-- Drop old settings table and create new one with hierarchy
DROP TABLE IF EXISTS settings;

CREATE TABLE settings (
    id BIGSERIAL PRIMARY KEY,

    -- Hierarchical scope (determines override precedence)
    scope TEXT NOT NULL CHECK (scope IN ('system', 'organization', 'project', 'user')),

    -- Foreign keys for scoped settings (nullable, depends on scope)
    organization_id BIGINT REFERENCES organization(id) ON DELETE CASCADE,
    project_id BIGINT REFERENCES project(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,

    -- Setting identification
    subsystem TEXT NOT NULL,  -- e.g., 'smtp', 'features', 'limits', 'branding'
    key TEXT NOT NULL,        -- e.g., 'rate_limit.api_calls', 'feature.webhooks'

    -- Value storage (JSONB for rich data types)
    value JSONB NOT NULL,

    -- Metadata
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE SET NULL
);

-- Unique index based on scope (supports expressions unlike UNIQUE constraint)
CREATE UNIQUE INDEX idx_settings_unique_scope ON settings(
    scope,
    COALESCE(organization_id, 0),
    COALESCE(project_id, 0),
    COALESCE(user_id, 0),
    subsystem,
    key
);

-- Indexes for fast lookups
CREATE INDEX idx_settings_scope ON settings(scope);
CREATE INDEX idx_settings_organization ON settings(organization_id) WHERE organization_id IS NOT NULL;
CREATE INDEX idx_settings_project ON settings(project_id) WHERE project_id IS NOT NULL;
CREATE INDEX idx_settings_user ON settings(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_settings_subsystem ON settings(subsystem);
CREATE INDEX idx_settings_key ON settings(key);

-- Check constraints to ensure scope consistency
ALTER TABLE settings ADD CONSTRAINT check_system_scope
    CHECK (scope != 'system' OR (organization_id IS NULL AND project_id IS NULL AND user_id IS NULL));

ALTER TABLE settings ADD CONSTRAINT check_organization_scope
    CHECK (scope != 'organization' OR (organization_id IS NOT NULL AND project_id IS NULL AND user_id IS NULL));

ALTER TABLE settings ADD CONSTRAINT check_project_scope
    CHECK (scope != 'project' OR (project_id IS NOT NULL AND user_id IS NULL));

ALTER TABLE settings ADD CONSTRAINT check_user_scope
    CHECK (scope != 'user' OR user_id IS NOT NULL);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER settings_updated_at
    BEFORE UPDATE ON settings
    FOR EACH ROW
    EXECUTE FUNCTION update_settings_updated_at();

-- Insert default system settings
INSERT INTO settings (scope, subsystem, key, value, description) VALUES
    ('system', 'features', 'webhooks.enabled', 'true', 'Enable webhooks feature'),
    ('system', 'features', 'api_keys.enabled', 'true', 'Enable API keys feature'),
    ('system', 'features', 'sso.enabled', 'false', 'Enable SSO authentication'),
    ('system', 'limits', 'rate_limit.api_calls', '{"limit": 1000, "window": "1h"}', 'API rate limit for authenticated requests'),
    ('system', 'limits', 'rate_limit.anonymous', '{"limit": 100, "window": "1h"}', 'API rate limit for anonymous requests'),
    ('system', 'limits', 'storage.max_mb', '1024', 'Maximum storage per project in MB'),
    ('system', 'limits', 'users.max_per_org', '100', 'Maximum users per organization'),
    ('system', 'branding', 'logo.url', '""', 'Default logo URL'),
    ('system', 'branding', 'primary_color', '"#3b82f6"', 'Default primary brand color'),
    ('system', 'notifications', 'email.enabled', 'true', 'Enable email notifications'),
    ('system', 'notifications', 'email.from_name', '"P402"', 'Default email from name');
