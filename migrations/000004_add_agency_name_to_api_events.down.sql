DROP TABLE IF EXISTS client_agencies;

ALTER TABLE api_events
DROP COLUMN IF EXISTS agency_name;
