-- name: GetSettingHierarchy :many
-- Get all matching settings in hierarchy order (system → org → project → user)
-- Returns ordered by precedence (higher scope overrides lower)
SELECT
    id,
    scope,
    organization_id,
    project_id,
    user_id,
    subsystem,
    key,
    value,
    description,
    created_at,
    updated_at,
    updated_by
FROM settings
WHERE subsystem = $1
    AND key = $2
    AND (
        scope = 'system'
        OR (scope = 'organization' AND organization_id = sqlc.narg('organization_id'))
        OR (scope = 'project' AND project_id = sqlc.narg('project_id'))
        OR (scope = 'user' AND user_id = sqlc.narg('user_id'))
    )
ORDER BY
    CASE scope
        WHEN 'user' THEN 4
        WHEN 'project' THEN 3
        WHEN 'organization' THEN 2
        WHEN 'system' THEN 1
    END DESC
LIMIT 1;

-- name: GetSettingsBySubsystem :many
-- Get all settings for a subsystem with hierarchy resolution
SELECT
    id,
    scope,
    organization_id,
    project_id,
    user_id,
    subsystem,
    key,
    value,
    description,
    created_at,
    updated_at,
    updated_by
FROM settings
WHERE subsystem = $1
    AND (
        scope = 'system'
        OR (scope = 'organization' AND organization_id = sqlc.narg('organization_id'))
        OR (scope = 'project' AND project_id = sqlc.narg('project_id'))
        OR (scope = 'user' AND user_id = sqlc.narg('user_id'))
    )
ORDER BY
    key,
    CASE scope
        WHEN 'user' THEN 4
        WHEN 'project' THEN 3
        WHEN 'organization' THEN 2
        WHEN 'system' THEN 1
    END DESC;

-- name: UpsertSystemSetting :one
-- Create or update a system-level setting
INSERT INTO settings (scope, subsystem, key, value, description, updated_by)
VALUES ('system', $1, $2, $3, sqlc.narg('description'), sqlc.narg('updated_by'))
ON CONFLICT (scope, COALESCE(organization_id, 0), COALESCE(project_id, 0), COALESCE(user_id, 0), subsystem, key)
DO UPDATE SET
    value = EXCLUDED.value,
    description = COALESCE(EXCLUDED.description, settings.description),
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: UpsertOrganizationSetting :one
-- Create or update an organization-level setting
INSERT INTO settings (scope, organization_id, subsystem, key, value, description, updated_by)
VALUES ('organization', $1, $2, $3, $4, sqlc.narg('description'), sqlc.narg('updated_by'))
ON CONFLICT (scope, COALESCE(organization_id, 0), COALESCE(project_id, 0), COALESCE(user_id, 0), subsystem, key)
DO UPDATE SET
    value = EXCLUDED.value,
    description = COALESCE(EXCLUDED.description, settings.description),
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: UpsertProjectSetting :one
-- Create or update a project-level setting
INSERT INTO settings (scope, project_id, subsystem, key, value, description, updated_by)
VALUES ('project', $1, $2, $3, $4, sqlc.narg('description'), sqlc.narg('updated_by'))
ON CONFLICT (scope, COALESCE(organization_id, 0), COALESCE(project_id, 0), COALESCE(user_id, 0), subsystem, key)
DO UPDATE SET
    value = EXCLUDED.value,
    description = COALESCE(EXCLUDED.description, settings.description),
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: UpsertUserSetting :one
-- Create or update a user-level setting
INSERT INTO settings (scope, user_id, subsystem, key, value, description, updated_by)
VALUES ('user', $1, $2, $3, $4, sqlc.narg('description'), sqlc.narg('updated_by'))
ON CONFLICT (scope, COALESCE(organization_id, 0), COALESCE(project_id, 0), COALESCE(user_id, 0), subsystem, key)
DO UPDATE SET
    value = EXCLUDED.value,
    description = COALESCE(EXCLUDED.description, settings.description),
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: DeleteSetting :exec
-- Delete a specific setting by ID
DELETE FROM settings WHERE id = $1;

-- name: ListAllSettings :many
-- List all settings (for admin interface)
SELECT
    id,
    scope,
    organization_id,
    project_id,
    user_id,
    subsystem,
    key,
    value,
    description,
    created_at,
    updated_at,
    updated_by
FROM settings
ORDER BY subsystem, key, scope;

-- name: ListSystemSettings :many
-- List all system-level settings
SELECT
    id,
    scope,
    organization_id,
    project_id,
    user_id,
    subsystem,
    key,
    value,
    description,
    created_at,
    updated_at,
    updated_by
FROM settings
WHERE scope = 'system'
ORDER BY subsystem, key;

-- name: ListOrganizationSettings :many
-- List settings for a specific organization (including system defaults)
SELECT
    id,
    scope,
    organization_id,
    project_id,
    user_id,
    subsystem,
    key,
    value,
    description,
    created_at,
    updated_at,
    updated_by
FROM settings
WHERE scope = 'system'
    OR (scope = 'organization' AND organization_id = $1)
ORDER BY subsystem, key, scope DESC;

-- name: ListProjectSettings :many
-- List settings for a specific project (including org and system defaults)
-- Note: Pass both project_id and org_id as parameters
SELECT
    id,
    scope,
    organization_id,
    project_id,
    user_id,
    subsystem,
    key,
    value,
    description,
    created_at,
    updated_at,
    updated_by
FROM settings
WHERE scope = 'system'
    OR (scope = 'organization' AND organization_id = $1)
    OR (scope = 'project' AND project_id = $2)
ORDER BY subsystem, key, scope DESC;

-- name: GetSettingByID :one
-- Get a specific setting by ID
SELECT
    id,
    scope,
    organization_id,
    project_id,
    user_id,
    subsystem,
    key,
    value,
    description,
    created_at,
    updated_at,
    updated_by
FROM settings
WHERE id = $1;
