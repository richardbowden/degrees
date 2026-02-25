-- name: GetCartBySessionToken :one
SELECT * FROM cart_sessions
WHERE session_token = $1
  AND expires_at > NOW();

-- name: GetCartByUserID :one
SELECT * FROM cart_sessions
WHERE user_id = $1
  AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateCartSession :one
INSERT INTO cart_sessions (user_id, session_token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddCartItem :one
INSERT INTO cart_items (cart_session_id, service_id, vehicle_id, quantity)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: RemoveCartItem :exec
DELETE FROM cart_items
WHERE id = $1;

-- name: ClearCart :exec
DELETE FROM cart_items
WHERE cart_session_id = $1;

-- name: ListCartItems :many
SELECT ci.id, ci.cart_session_id, ci.service_id, ci.vehicle_id,
       ci.quantity, ci.created_at,
       s.name AS service_name,
       COALESCE(spt.price, s.base_price) AS service_price
FROM cart_items ci
JOIN services s ON s.id = ci.service_id
LEFT JOIN vehicles v ON v.id = ci.vehicle_id
LEFT JOIN service_price_tiers spt ON spt.service_id = s.id AND spt.vehicle_category_id = v.vehicle_category_id
WHERE ci.cart_session_id = $1
ORDER BY ci.created_at;

-- name: AddCartItemOption :one
INSERT INTO cart_item_options (cart_item_id, service_option_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveCartItemOption :exec
DELETE FROM cart_item_options
WHERE cart_item_id = $1 AND service_option_id = $2;
