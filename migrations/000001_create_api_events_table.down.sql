-- 000001_create_api_events_table.down.sql
DROP INDEX IF EXISTS idx_api_events_client_id;
DROP INDEX IF EXISTS idx_api_events_timestamp;
DROP TABLE IF EXISTS api_events;
