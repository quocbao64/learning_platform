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

// @Summary Đăng ký khóa học
// @Tags Enrollments
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "ID khoá học"
// @Success 201 {object} object{message=string} "Đăng ký học thành công"
// @Failure 400 {object} response.ErrorResponse "Khóa học đã đầy hoặc đã đăng ký trước đó"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /courses/{course_id}/enroll [post]
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

// @Summary Danh sách khóa học đã đăng ký
// @Tags Enrollments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{enrollments=[]models.Enrollment} "Danh sách đăng ký học"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /me/enrollments [get]
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

// @Summary Hủy đăng ký học
// @Tags Enrollments
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "ID khóa học"
// @Success 200 {object} object{message=string} "Hủy đăng ký thành công"
// @Failure 400 {object} response.ErrorResponse "Không tìm thấy đăng ký học"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /courses/{course_id}/enroll [delete]
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
