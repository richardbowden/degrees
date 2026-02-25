-- Vehicle categories (e.g. Sedan/Hatchback, SUV/Wagon, 4WD) for size-based pricing

CREATE TABLE vehicle_categories (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
SELECT add_updated_at_trigger('vehicle_categories');

CREATE TABLE service_price_tiers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    vehicle_category_id BIGINT NOT NULL REFERENCES vehicle_categories(id) ON DELETE CASCADE,
    price BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(service_id, vehicle_category_id)
);
CREATE INDEX idx_spt_service ON service_price_tiers(service_id);
CREATE INDEX idx_spt_category ON service_price_tiers(vehicle_category_id);

ALTER TABLE vehicles ADD COLUMN vehicle_category_id BIGINT REFERENCES vehicle_categories(id);
