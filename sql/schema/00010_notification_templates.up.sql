-- Templates table
CREATE TABLE if not exists template (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    scope_type TEXT NOT NULL, ---System', 'Organization', 'Project')
    version INTEGER NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    created_by BIGINT references users(id),
    updated_by BIGINT references users(id)
);

SELECT add_updated_at_trigger('template');
