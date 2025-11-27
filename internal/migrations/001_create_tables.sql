-- +goose Up
CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    url TEXT NOT NULL,
    interval_seconds INTEGER NOT NULL CHECK (interval_seconds >= 10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE check_results (
    id SERIAL PRIMARY KEY,
    service_id INTEGER REFERENCES services(id) ON DELETE CASCADE,
    status_code INTEGER,
    response_time_ms INTEGER NOT NULL,
    checked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_services_name ON services(name);
CREATE INDEX idx_check_results_service_id ON check_results(service_id);
CREATE INDEX idx_check_results_checked_at ON check_results(checked_at);

-- +goose Down
DROP TABLE IF EXISTS check_results;
DROP TABLE IF EXISTS services;