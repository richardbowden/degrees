-- Migration F: Service History

CREATE TABLE IF NOT EXISTS service_records (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id),
    customer_id BIGINT NOT NULL REFERENCES customer_profiles(id),
    vehicle_id BIGINT REFERENCES vehicles(id),
    completed_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_service_records_booking_id ON service_records(booking_id);
CREATE INDEX idx_service_records_customer_id ON service_records(customer_id);
CREATE INDEX idx_service_records_vehicle_id ON service_records(vehicle_id);
SELECT add_updated_at_trigger('service_records');

CREATE TABLE IF NOT EXISTS service_notes (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_record_id BIGINT NOT NULL REFERENCES service_records(id) ON DELETE CASCADE,
    note_type TEXT NOT NULL,
    content TEXT NOT NULL,
    is_visible_to_customer BOOLEAN NOT NULL DEFAULT true,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_service_notes_service_record_id ON service_notes(service_record_id);

CREATE TABLE IF NOT EXISTS service_products_used (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_record_id BIGINT NOT NULL REFERENCES service_records(id) ON DELETE CASCADE,
    product_name TEXT NOT NULL,
    notes TEXT
);

CREATE INDEX idx_service_products_used_service_record_id ON service_products_used(service_record_id);

CREATE TABLE IF NOT EXISTS service_photos (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_record_id BIGINT NOT NULL REFERENCES service_records(id) ON DELETE CASCADE,
    photo_type TEXT NOT NULL,
    url TEXT NOT NULL,
    caption TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_service_photos_service_record_id ON service_photos(service_record_id);
