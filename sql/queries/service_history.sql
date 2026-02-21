-- name: CreateServiceRecord :one
INSERT INTO service_records (booking_id, customer_id, vehicle_id, completed_date)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetServiceRecordByID :one
SELECT * FROM service_records
WHERE id = $1;

-- name: ListServiceRecordsByCustomer :many
SELECT * FROM service_records
WHERE customer_id = $1
ORDER BY completed_date DESC;

-- name: ListServiceRecordsByBooking :many
SELECT * FROM service_records
WHERE booking_id = $1
ORDER BY completed_date DESC;

-- name: CreateServiceNote :one
INSERT INTO service_notes (service_record_id, note_type, content, is_visible_to_customer, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListServiceNotes :many
SELECT * FROM service_notes
WHERE service_record_id = $1
ORDER BY created_at;

-- name: CreateServiceProductUsed :one
INSERT INTO service_products_used (service_record_id, product_name, notes)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListServiceProductsUsed :many
SELECT * FROM service_products_used
WHERE service_record_id = $1;

-- name: CreateServicePhoto :one
INSERT INTO service_photos (service_record_id, photo_type, url, caption)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListServicePhotos :many
SELECT * FROM service_photos
WHERE service_record_id = $1
ORDER BY created_at;
