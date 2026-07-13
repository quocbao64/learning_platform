package models

import "net/http"

var (
	ErrCourseAlreadyExists  = NewAppError("COURSE_ALREADY_EXISTS", "course already exists", http.StatusBadRequest)
	ErrCourseNotFound       = NewAppError("COURSE_NOT_FOUND", "course not found", http.StatusBadRequest)
	ErrCourseNotPublished   = NewAppError("COURSE_NOT_PUBLISHED", "course not published", http.StatusBadRequest)
	ErrCourseFull           = NewAppError("COURSE_FULL", "course full", http.StatusBadRequest)
	ErrCourseInvalidStatus  = NewAppError("COURSE_INVALID_STATUS", "course invalid status", http.StatusBadRequest)
	ErrCourseHasEnrollments = NewAppError("COURSE_HAS_ENROLLMENTS", "course has enrollments", http.StatusBadRequest)

	ErrEnrollmentAlreadyExists = NewAppError("ENROLLMENT_ALREADY_EXISTS", "enrollment already exists", http.StatusBadRequest)
	ErrEnrollmentNotFound      = NewAppError("ENROLLMENT_NOT_FOUND", "enrollment not found", http.StatusBadRequest)

	ErrUserNotFound       = NewAppError("USER_NOT_FOUND", "user not found", http.StatusBadRequest)
	ErrEmailAlreadyExists = NewAppError("EMAIL_ALREADY_EXISTS", "email already exists", http.StatusBadRequest)
	ErrInvalidRole        = NewAppError("INVALID_ROLE", "invalid role", http.StatusBadRequest)

	ErrCacheMiss = NewAppError("CACHE_MISS", "cache miss", http.StatusBadRequest)

	ErrInvalidCredentials = NewAppError("INVALID_CREDENTIALS", "invalid credentials", http.StatusUnauthorized)

	ErrInternal = NewAppError("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)

	ErrForbidden = NewAppError("FORBIDDEN", "forbidden", http.StatusForbidden)
)
