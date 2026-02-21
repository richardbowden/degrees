DROP TABLE IF EXISTS cart_item_options;

DROP INDEX IF EXISTS idx_cart_items_cart_session_id;
DROP TABLE IF EXISTS cart_items;

DROP INDEX IF EXISTS idx_cart_sessions_expires_at;
DROP INDEX IF EXISTS idx_cart_sessions_user_id;
DROP INDEX IF EXISTS idx_cart_sessions_token;
DROP TABLE IF EXISTS cart_sessions;
