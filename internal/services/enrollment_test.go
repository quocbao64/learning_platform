package services

import (
	"context"
	"errors"
	"learning-platform/internal/models"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

func TestCreateEnrollment(t *testing.T) {
	tests := []struct {
		name      string
		course    *models.Course
		courseErr error
		existing  *models.Enrollment
		findErr   error
		courseID  int64
		seatsOK   bool
		createErr error
		userID    int64
		expectErr error
	}{
		{
			name: "success",
			course: &models.Course{
				ID:     1,
				Status: models.CourseStatusPublished,
			},
			courseID:  1,
			seatsOK:   true,
			userID:    1,
			expectErr: nil,
		},
		{
			name:      "course not found",
			course:    nil,
			courseErr: models.ErrCourseNotFound,
			courseID:  99,
			userID:    1,
			expectErr: models.ErrCourseNotFound,
		},
		{
			name: "course not published",
			course: &models.Course{
				ID:     1,
				Status: models.CourseStatusDraft,
			},
			courseID:  1,
			userID:    1,
			expectErr: models.ErrCourseNotPublished,
		},
		{
			name: "already enrolled",
			course: &models.Course{
				ID:     1,
				Status: models.CourseStatusPublished,
			},
			existing:  &models.Enrollment{ID: 5},
			courseID:  1,
			userID:    1,
			expectErr: models.ErrEnrollmentAlreadyExists,
		},
		{
			name: "course full",
			course: &models.Course{
				ID:     1,
				Status: models.CourseStatusPublished,
			},
			courseID:  1,
			seatsOK:   false,
			userID:    1,
			expectErr: models.ErrCourseFull,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseRepo := &MockCourseRepository{}
			enrollRepo := &MockEnrollmentRepository{}
			txManager := &MockTxManager{}
			cache := &MockCache{}

			courseRepo.On("GetByID", mock.Anything, tt.courseID).
				Return(tt.course, tt.courseErr)

			if tt.course != nil && tt.course.Status == models.CourseStatusPublished {
				enrollRepo.On("FindByUserIDAndCourseID", mock.Anything, tt.userID, tt.courseID).
					Return(tt.existing, tt.findErr)

				if tt.existing == nil {
					txManager.On("ExecTx", mock.Anything, mock.Anything).
						Return(func(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
							return fn(ctx, nil)
						})

					courseRepo.On("DecrementSeatsTx", mock.Anything, mock.Anything, tt.courseID).
						Return(tt.seatsOK, nil)

					if tt.seatsOK {
						enrollRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
							Return(&models.Enrollment{ID: 1}, tt.createErr)
					}
				}
			}

			svc := NewEnrollmentService(enrollRepo, courseRepo, txManager, cache)

			_, err := svc.CreateEnrollment(context.Background(), &models.Enrollment{
				UserID:   tt.userID,
				CourseID: tt.courseID,
				Status:   models.EnrollmentStatusActive,
			})

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("CreateEnrollment() error = %v, want %v", err, tt.expectErr)
			}
		})
	}
}
