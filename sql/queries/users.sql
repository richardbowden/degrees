-- name: CreateUser :one
INSERT INTO accounts (acc_num, first_name, middle_name, surname, email, password_hash, acc_type)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- -- name: CreateOrganisation :one
-- INSERT INTO organisations (display_name, slug, owner_id)
-- VALUES ($1, $2, $3)
-- RETURNING *;
--
-- -- name: GetUser :one
-- SELECT * from users where email = $1;
--
-- -- name: LinkProjectToOrganisation :exec
-- INSERT INTO organisation_projects (org_id, project_id) VALUES ($1, $2);
--
-- -- name: AddEmailAddress :exec
-- INSERT INTO email_addresses (user_id, email) VALUES ($1, $2);
--
-- -- name: CheckEmailExists :one
-- SELECT exists(select 1 from users where email = $1);
--
-- -- name: AddUserValidationToken :exec
-- INSERT INTO user_validation_tokens (user_id, token, expires_at) VALUES ($1, $2, $3);
--
-- -- name: DeleteUserValidationToken :exec
-- DELETE FROM user_validation_tokens where token = $1 and user_id = $2;
--
-- -- name: SetEmailVerifiedForUser :exec
-- -- UPDATE users set verified = $1 where id = $2;
--
-- -- name: SetUserPendingEmailValidation :exec
-- UPDATE users set sign_up_stage = 1 where id = $1;
