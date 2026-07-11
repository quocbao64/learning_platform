package models

import "time"

type Progress struct {
	ID           int64     `json:"id" db:"id"`
	EnrollmentID int64     `json:"enrollment_id" db:"enrollment_id"`
	LessonID     int64     `json:"lesson_id" db:"lesson_id"`
	Status       string    `json:"status" db:"status"`
	CompletedAt  time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

func (p Progress) TableName() string {
	return "progress"
}

var (
	ProgressStatusCompleted  = "completed"
	ProgressStatusNotStarted = "not_started"
)

type CourseProgress struct {
	EnrollmentID     int64   `json:"enrollment_id" db:"enrollment_id"`
	TotalLessons     int     `json:"total_lessons" db:"total_lessons"`
	CompletedCount   int     `json:"completed_count" db:"completed_count"`
	PercentCompleted float64 `json:"percent_completed" db:"percent_completed"`
}
