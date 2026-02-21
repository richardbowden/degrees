-- name: CreateBooking :one
INSERT INTO bookings (
    customer_id, vehicle_id, scheduled_date, scheduled_time,
    estimated_duration_mins, status, payment_status,
    subtotal, deposit_amount, total_amount,
    stripe_payment_intent_id, stripe_deposit_intent_id, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetBookingByID :one
SELECT b.*,
       cp.user_id AS customer_user_id,
       cp.phone AS customer_phone,
       v.make AS vehicle_make,
       v.model AS vehicle_model,
       v.rego AS vehicle_rego
FROM bookings b
JOIN customer_profiles cp ON cp.id = b.customer_id
LEFT JOIN vehicles v ON v.id = b.vehicle_id
WHERE b.id = $1;

-- name: ListBookingsByCustomer :many
SELECT * FROM bookings
WHERE customer_id = $1
ORDER BY scheduled_date DESC, scheduled_time DESC;

-- name: ListBookingsByDateRange :many
SELECT * FROM bookings
WHERE scheduled_date >= $1 AND scheduled_date <= $2
ORDER BY scheduled_date, scheduled_time;

-- name: ListBookingsForDate :many
SELECT * FROM bookings
WHERE scheduled_date = $1
  AND status NOT IN ('cancelled')
ORDER BY scheduled_time;

-- name: UpdateBookingStatus :one
UPDATE bookings
SET status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateBookingPaymentStatus :one
UPDATE bookings
SET payment_status = $2
WHERE id = $1
RETURNING *;

-- name: CreateBookingService :one
INSERT INTO booking_services (booking_id, service_id, price_at_booking)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateBookingServiceOption :one
INSERT INTO booking_service_options (booking_service_id, service_option_id, price_at_booking)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListBookingServices :many
SELECT bs.id, bs.booking_id, bs.service_id, bs.price_at_booking,
       s.name AS service_name, s.slug AS service_slug
FROM booking_services bs
JOIN services s ON s.id = bs.service_id
WHERE bs.booking_id = $1;

-- name: ListBookingServiceOptions :many
SELECT bso.id, bso.booking_service_id, bso.service_option_id,
       bso.price_at_booking,
       so.name AS option_name
FROM booking_service_options bso
JOIN service_options so ON so.id = bso.service_option_id
WHERE bso.booking_service_id = $1;
