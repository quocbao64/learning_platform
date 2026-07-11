CREATE TABLE progress (
    id BIGSERIAL PRIMARY KEY,
    enrollment_id BIGSERIAL NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    lesson_id BIGSERIAL NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    status VARCHAR(10) NOT NULL DEFAULT 'not_started',
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_progress_enrollment_id ON progress(enrollment_id);
CREATE INDEX idx_progress_lesson_id ON progress(lesson_id);
CREATE UNIQUE INDEX idx_progress_enrollment_lesson ON progress(enrollment_id, lesson_id);