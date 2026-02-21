-- Migration E: Schedule

CREATE TABLE IF NOT EXISTS schedule_config (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    is_open BOOLEAN NOT NULL DEFAULT true,
    buffer_minutes INT NOT NULL DEFAULT 30,
    UNIQUE(day_of_week)
);

CREATE TABLE IF NOT EXISTS schedule_blackouts (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    date DATE NOT NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_schedule_blackouts_date ON schedule_blackouts(date);
