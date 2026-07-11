package models

import (
	"errors"
	"time"
)

type Course struct {
	ID           int64      `db:"id" json:"id"`
	InstructorID int64      `db:"instructor_id" json:"instructor_id"`
	Title        string     `db:"title" json:"title"`
	Description  string     `db:"description" json:"description"`
	Status       string     `db:"status" json:"status"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at"`

	Lessons []*Lesson `db:"-" json:"lessons"`
}

func (course *Course) TableName() string {
	return "courses"
}

var (
	ErrCourseAlreadyExists = errors.New("course already exists")
	ErrCourseNotFound      = errors.New("course not found")
)

var (
	CourseStatusDraft     = "draft"
	CourseStatusPublished = "published"
	CourseStatusArchived  = "archived"
)
