-- name: CreateVerificationToken :exec
INSERT INTO verification (
    token,
    user_id,
    created_at)
VALUES ($1, $2, $3);

-- name: GetToken :one
SELECT * from verification where token = $1;

-- name: DeleteToken :exec
DELETE from verification where token = $1;
