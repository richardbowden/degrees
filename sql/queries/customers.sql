-- name: GetCustomerProfileByUserID :one
SELECT * FROM customer_profiles
WHERE user_id = $1;

-- name: CreateCustomerProfile :one
INSERT INTO customer_profiles (user_id, phone, address, suburb, postcode, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateCustomerProfile :one
UPDATE customer_profiles
SET phone = $2, address = $3, suburb = $4, postcode = $5, notes = $6
WHERE id = $1
RETURNING *;

-- name: ListCustomers :many
SELECT * FROM customer_profiles
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListVehiclesByCustomer :many
SELECT * FROM vehicles
WHERE customer_id = $1
ORDER BY is_primary DESC, created_at DESC;

-- name: GetVehicleByID :one
SELECT * FROM vehicles
WHERE id = $1;

-- name: CreateVehicle :one
INSERT INTO vehicles (
    customer_id, make, model, year, colour, rego,
    paint_type, condition_notes, is_primary, vehicle_category_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateVehicle :one
UPDATE vehicles
SET make = $2, model = $3, year = $4, colour = $5, rego = $6,
    paint_type = $7, condition_notes = $8, is_primary = $9,
    vehicle_category_id = $10
WHERE id = $1
RETURNING *;

-- name: DeleteVehicle :exec
DELETE FROM vehicles
WHERE id = $1;
