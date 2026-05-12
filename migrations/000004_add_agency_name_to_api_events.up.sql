ALTER TABLE api_events
ADD COLUMN IF NOT EXISTS agency_name TEXT;

CREATE TABLE IF NOT EXISTS client_agencies (
    client_id TEXT PRIMARY KEY,
    agency_name TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
