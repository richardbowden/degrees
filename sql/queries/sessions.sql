-- name: CreateSession :one
INSERT INTO sessions (
    user_id,
    session_token,
    expires_at,
    user_agent,
    ip_address
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE session_token = $1
  AND expires_at > NOW()
LIMIT 1;

-- name: UpdateSessionActivity :exec
UPDATE sessions
SET last_activity_at = NOW()
WHERE session_token = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_token = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < NOW();

-- name: GetUserActiveSessions :many
SELECT * FROM sessions
WHERE user_id = $1
  AND expires_at > NOW()
ORDER BY last_activity_at DESC;
