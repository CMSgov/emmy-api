CREATE TABLE IF NOT EXISTS batch_student_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_student_id UUID NOT NULL REFERENCES batch_students(id) ON DELETE CASCADE,
    results JSONB,
    found_enrollment BOOLEAN DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_batch_student_results_batch_student_id ON batch_student_results(batch_student_id);
