
-- Organizations table
CREATE TABLE if not exists organization (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(39) NOT NULL UNIQUE,
    description TEXT NULL,
    avatar VARCHAR(500) NULL,
    website VARCHAR(500) NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

SELECT add_updated_at_trigger('organization');
CREATE INDEX if not exists org_idx_slug ON organization (slug);

CREATE TABLE if not exists org_user_membership (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id BIGINT NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    -- role needs more thinking and work
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'member', 'viewer')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (user_id, organization_id)
);


