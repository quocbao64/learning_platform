package models

import (
	"errors"
	"time"
)

type Enrollment struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	CourseID   int64     `json:"course_id" db:"course_id"`
	Status     string    `json:"status" db:"status"`
	EnrolledAt time.Time `json:"enrolled_at" db:"enrolled_at"`
}

func (e Enrollment) TableName() string {
	return "enrollments"
}

var (
	ErrEnrollmentAlreadyExists = errors.New("enrollment already exists")
)

var (
	EnrollmentStatusActive    = "active"
	EnrollmentStatusCompleted = "completed"
	EnrollmentStatusCancelled = "cancelled"
)
