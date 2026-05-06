CREATE TABLE IF NOT EXISTS enrollment_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id TEXT NOT NULL UNIQUE,
    submitted_by TEXT NOT NULL,
    callback_url TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS batch_students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_db_id UUID NOT NULL REFERENCES enrollment_batches(id) ON DELETE CASCADE,
    record_id TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth TEXT NOT NULL,
    ssn TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_batch_students_batch_db_id ON batch_students(batch_db_id);
