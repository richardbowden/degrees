-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (
    token,
    user_id,
    expires_at,
    created_at)
VALUES ($1, $2, $3, $4);

-- name: GetPasswordResetToken :one
SELECT * FROM password_reset_tokens WHERE token = $1 AND expires_at > NOW();

-- name: DeletePasswordResetToken :exec
DELETE FROM password_reset_tokens WHERE token = $1;

-- name: DeleteExpiredPasswordResetTokens :exec
DELETE FROM password_reset_tokens WHERE expires_at < NOW();

-- name: DeleteUserPasswordResetTokens :exec
DELETE FROM password_reset_tokens WHERE user_id = $1;
