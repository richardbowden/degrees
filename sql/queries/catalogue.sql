-- name: ListCategories :many
SELECT * FROM service_categories
ORDER BY sort_order, name;

-- name: GetCategoryBySlug :one
SELECT * FROM service_categories
WHERE slug = $1;

-- name: CreateCategory :one
INSERT INTO service_categories (name, slug, description, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCategory :one
UPDATE service_categories
SET name = $2, slug = $3, description = $4, sort_order = $5
WHERE id = $1
RETURNING *;

-- name: ListServices :many
SELECT * FROM services
WHERE is_active = true
ORDER BY sort_order, name;

-- name: ListServicesByCategory :many
SELECT * FROM services
WHERE category_id = $1 AND is_active = true
ORDER BY sort_order, name;

-- name: GetServiceBySlug :one
SELECT s.*, sc.name AS category_name
FROM services s
JOIN service_categories sc ON sc.id = s.category_id
WHERE s.slug = $1;

-- name: GetServiceByID :one
SELECT * FROM services
WHERE id = $1;

-- name: CreateService :one
INSERT INTO services (
    category_id, name, slug, description, short_desc,
    base_price, duration_minutes, is_active, sort_order
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateService :one
UPDATE services
SET category_id = $2, name = $3, slug = $4, description = $5,
    short_desc = $6, base_price = $7, duration_minutes = $8,
    is_active = $9, sort_order = $10
WHERE id = $1
RETURNING *;

-- name: DeleteService :one
UPDATE services
SET is_active = false
WHERE id = $1
RETURNING *;

-- name: ListServiceOptions :many
SELECT * FROM service_options
WHERE service_id = $1 AND is_active = true
ORDER BY sort_order, name;

-- name: CreateServiceOption :one
INSERT INTO service_options (service_id, name, description, price, is_active, sort_order)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateServiceOption :one
UPDATE service_options
SET name = $2, description = $3, price = $4, is_active = $5, sort_order = $6
WHERE id = $1
RETURNING *;

-- name: DeleteServiceOption :one
UPDATE service_options
SET is_active = false
WHERE id = $1
RETURNING *;
