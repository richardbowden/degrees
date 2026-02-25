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

-- ========================================
-- Vehicle Categories
-- ========================================

-- name: ListVehicleCategories :many
SELECT * FROM vehicle_categories
ORDER BY sort_order, name;

-- name: GetVehicleCategoryByID :one
SELECT * FROM vehicle_categories
WHERE id = $1;

-- name: CreateVehicleCategory :one
INSERT INTO vehicle_categories (name, slug, description, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateVehicleCategory :one
UPDATE vehicle_categories
SET name = $2, slug = $3, description = $4, sort_order = $5
WHERE id = $1
RETURNING *;

-- name: DeleteVehicleCategory :exec
DELETE FROM vehicle_categories
WHERE id = $1;

-- ========================================
-- Service Price Tiers
-- ========================================

-- name: ListPriceTiersByService :many
SELECT spt.id, spt.service_id, spt.vehicle_category_id, spt.price, spt.created_at,
       vc.name AS category_name, vc.slug AS category_slug
FROM service_price_tiers spt
JOIN vehicle_categories vc ON vc.id = spt.vehicle_category_id
WHERE spt.service_id = $1
ORDER BY vc.sort_order, vc.name;

-- name: UpsertPriceTier :one
INSERT INTO service_price_tiers (service_id, vehicle_category_id, price)
VALUES ($1, $2, $3)
ON CONFLICT (service_id, vehicle_category_id)
DO UPDATE SET price = EXCLUDED.price
RETURNING *;

-- name: DeletePriceTier :exec
DELETE FROM service_price_tiers
WHERE service_id = $1 AND vehicle_category_id = $2;

-- name: GetPriceTier :one
SELECT spt.id, spt.service_id, spt.vehicle_category_id, spt.price, spt.created_at,
       vc.name AS category_name, vc.slug AS category_slug
FROM service_price_tiers spt
JOIN vehicle_categories vc ON vc.id = spt.vehicle_category_id
WHERE spt.service_id = $1 AND spt.vehicle_category_id = $2;

-- name: DeletePriceTiersByService :exec
DELETE FROM service_price_tiers
WHERE service_id = $1;
