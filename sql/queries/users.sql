-- First create the email, then create the user in separate operations
-- name: CreateUserEmail :one
INSERT INTO user_email (
    user_id,
    email,
    enabled,
    is_verified
    ) VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: CreateUser :one
INSERT INTO users (
    first_name,
    middle_name,
    surname,
    username,
    login_email,
    primary_email_id,
    password_hash,
    sign_up_stage
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING *;

-- name: UserExists :one
SELECT
    EXISTS(SELECT 1 FROM users WHERE users.login_email = $1) AS email_exists,
    EXISTS(SELECT 1 FROM users WHERE users.username = $2) AS username_exists;

-- name: EmailExists :one
SELECT
    EXISTS(select 1 from users where users.login_email = $1);

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE login_email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = $2,
    middle_name = $3,
    surname = $4,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: UpdateUserEnabled :one
UPDATE users
SET enabled = $2,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: UpdateUserSignUpStage :one
UPDATE users
SET sign_up_stage = $2,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
set password_hash = $2
where id = $1;

-- name: UpdateUserSysop :one
UPDATE users
SET sysop = $2,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: IsFirstUser :one
SELECT (COUNT(*) = 1) AS is_first_user FROM users;

-- name: ListAllUsers :many
SELECT * FROM users
ORDER BY created_on DESC;
