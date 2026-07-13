package handlers

import (
	"learning-platform/internal/handlers/response"
	"learning-platform/internal/middleware"
	"learning-platform/internal/models"
	"learning-platform/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	enrollmentService services.EnrollmentService
}

func NewEnrollmentHandler(enrollmentService services.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentService: enrollmentService,
	}
}

func (h *EnrollmentHandler) RegisterRoutes(router *gin.RouterGroup, authMW gin.HandlerFunc) {
	router.GET("/me/enrollments", authMW, h.listEnrollments)
	router.POST("/courses/:course_id/enroll", authMW, h.enroll)
	router.DELETE("/courses/:course_id/enroll", authMW, h.cancel)
}

func (h *EnrollmentHandler) enroll(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	userID := middleware.UserID(c)

	_, err = h.enrollmentService.CreateEnrollment(c.Request.Context(), &models.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		Status:     models.EnrollmentStatusActive,
		EnrolledAt: time.Now(),
	})

	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Enrollment successfully",
	})
}

func (h *EnrollmentHandler) listEnrollments(c *gin.Context) {
	userID := middleware.UserID(c)

	enrollments, err := h.enrollmentService.ListEnrollment(c.Request.Context(), &services.EnrollmentFilter{
		UserID: userID,
	})

	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments": enrollments,
	})
}

func (h *EnrollmentHandler) cancel(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	userID := middleware.UserID(c)

	err = h.enrollmentService.DeleteEnrollment(c.Request.Context(), userID, courseID)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Enrollment canceled successfully",
	})
}
