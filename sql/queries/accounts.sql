-- name: CreateAccount :one
INSERT INTO accounts (account_number, first_name, middle_name, surname, email, password_hash, acc_type)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: AddAccountEmailAddress :exec
INSERT INTO email_addresses (email, account_id) 
VALUES ($1, $2);

-- name: CheckAccountActiveEmailExists :one
SELECT exists(select 1 from accounts where email = $1);


-- name: AddAccountEmailValidationToken :exec
INSERT INTO email_verification_code (account_id, code, expires_at) VALUES ($1, $2, $3);
--
-- -- name: DeleteUserValidationToken :exec
-- DELETE FROM user_validation_tokens where token = $1 and user_id = $2;
--
-- -- name: SetEmailVerifiedForUser :exec
-- -- UPDATE users set verified = $1 where id = $2;
--
-- -- name: SetUserPendingEmailValidation :exec
-- UPDATE users set sign_up_stage = 1 where id = $1;





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
--
--
