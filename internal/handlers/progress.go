package handlers

import (
	"learning-platform/internal/handlers/response"
	"learning-platform/internal/middleware"
	"learning-platform/internal/models"
	"learning-platform/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	progressService services.ProgressService
}

func NewProgressHandler(progressService services.ProgressService) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
	}
}

func (h *ProgressHandler) RegisterRoutes(rg *gin.RouterGroup, authMW gin.HandlerFunc) {
	rg.PATCH("/enrollments/:enrollment_id/lessons/:lesson_id/progress", authMW, h.updateProgress)
	rg.GET("/enrollments/:enrollment_id/progress", authMW, h.getProgress)
}

// @Summary Cập nhật tiến độ khóa học
// @Tags Progress
// @Produce json
// @Security BearerAuth
// @Param enrollment_id path int true "ID đăng ký khóa học"
// @Param lesson_id path int true "ID bài học"
// @Success 200 {object} object{message=string} "Cập nhật tiến độ thành công"
// @Failure 400 {object} response.ErrorResponse "ID đăng ký khóa học hoặc ID bài học không hợp lệ"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /enrollments/{enrollment_id}/lessons/{lesson_id}/progress [patch]
func (h *ProgressHandler) updateProgress(c *gin.Context) {
	enrollmentId, err := strconv.ParseInt(c.Param("enrollment_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	lessonId, err := strconv.ParseInt(c.Param("lesson_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	_, err = h.progressService.UpdateProgress(c, &models.Progress{
		EnrollmentID: enrollmentId,
		LessonID:     lessonId,
		Status:       models.ProgressStatusCompleted,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress updated successfully",
	})
}

// @Summary Xem tiến độ học tập
// @Tags Progress
// @Produce json
// @Security BearerAuth
// @Param enrollment_id path int true "ID đăng ký khóa học"
// @Success 200 {object} object{progress=models.CourseProgress} "Tiến độ học tập"
// @Failure 400 {object} response.ErrorResponse "ID đăng ký khóa học không hợp lệ"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /enrollments/{enrollment_id}/progress [get]
func (h *ProgressHandler) getProgress(c *gin.Context) {
	enrollmentId, err := strconv.ParseInt(c.Param("enrollment_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	userID := middleware.UserID(c)

	progress, err := h.progressService.GetProgress(c, userID, enrollmentId)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"progress": progress,
	})
}
