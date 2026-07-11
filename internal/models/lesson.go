package models

import "time"

type Lesson struct {
	ID         int64     `json:"id" db:"id"`
	CourseID   int64     `json:"course_id" db:"course_id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	OrderIndex int64     `json:"order_index" db:"order_index"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func (l Lesson) TableName() string {
	return "lessons"
}
