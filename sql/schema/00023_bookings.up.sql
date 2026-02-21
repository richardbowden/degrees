-- Migration D: Bookings

CREATE TYPE booking_status AS ENUM (
    'pending_payment',
    'deposit_paid',
    'confirmed',
    'in_progress',
    'completed',
    'cancelled'
);

CREATE TYPE payment_status AS ENUM (
    'pending',
    'deposit_paid',
    'fully_paid',
    'refunded',
    'partially_refunded'
);

CREATE TABLE IF NOT EXISTS bookings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customer_profiles(id),
    vehicle_id BIGINT REFERENCES vehicles(id),
    scheduled_date DATE NOT NULL,
    scheduled_time TIME NOT NULL,
    estimated_duration_mins INT NOT NULL,
    status booking_status NOT NULL DEFAULT 'pending_payment',
    payment_status payment_status NOT NULL DEFAULT 'pending',
    subtotal BIGINT NOT NULL,
    deposit_amount BIGINT NOT NULL,
    total_amount BIGINT NOT NULL,
    stripe_payment_intent_id TEXT,
    stripe_deposit_intent_id TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bookings_customer_id ON bookings(customer_id);
CREATE INDEX idx_bookings_vehicle_id ON bookings(vehicle_id);
CREATE INDEX idx_bookings_scheduled_date ON bookings(scheduled_date);
CREATE INDEX idx_bookings_status ON bookings(status);
SELECT add_updated_at_trigger('bookings');

CREATE TABLE IF NOT EXISTS booking_services (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    service_id BIGINT NOT NULL REFERENCES services(id),
    price_at_booking BIGINT NOT NULL,
    UNIQUE(booking_id, service_id)
);

CREATE INDEX idx_booking_services_booking_id ON booking_services(booking_id);

CREATE TABLE IF NOT EXISTS booking_service_options (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    booking_service_id BIGINT NOT NULL REFERENCES booking_services(id) ON DELETE CASCADE,
    service_option_id BIGINT NOT NULL REFERENCES service_options(id),
    price_at_booking BIGINT NOT NULL
);

CREATE INDEX idx_booking_service_options_booking_service_id ON booking_service_options(booking_service_id);
