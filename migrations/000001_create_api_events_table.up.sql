-- 000001_create_api_events_table.up.sql
CREATE TABLE IF NOT EXISTS api_events (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    endpoint TEXT NOT NULL,
    data_source TEXT NOT NULL,
    client_id TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    success BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_events_timestamp ON api_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_api_events_client_id ON api_events(client_id);
