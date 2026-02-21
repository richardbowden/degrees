-- Migration C: Cart

CREATE TABLE IF NOT EXISTS cart_sessions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    session_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cart_sessions_token ON cart_sessions(session_token);
CREATE INDEX idx_cart_sessions_user_id ON cart_sessions(user_id);
CREATE INDEX idx_cart_sessions_expires_at ON cart_sessions(expires_at);

CREATE TABLE IF NOT EXISTS cart_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cart_session_id BIGINT NOT NULL REFERENCES cart_sessions(id) ON DELETE CASCADE,
    service_id BIGINT NOT NULL REFERENCES services(id),
    vehicle_id BIGINT REFERENCES vehicles(id),
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cart_items_cart_session_id ON cart_items(cart_session_id);

CREATE TABLE IF NOT EXISTS cart_item_options (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cart_item_id BIGINT NOT NULL REFERENCES cart_items(id) ON DELETE CASCADE,
    service_option_id BIGINT NOT NULL REFERENCES service_options(id),
    UNIQUE(cart_item_id, service_option_id)
);
