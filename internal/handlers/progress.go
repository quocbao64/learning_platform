package handlers

import (
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

func (h *ProgressHandler) updateProgress(c *gin.Context) {
	enrollmentId, err := strconv.ParseInt(c.Param("enrollment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lessonId, err := strconv.ParseInt(c.Param("lesson_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.progressService.UpdateProgress(c, &models.Progress{
		EnrollmentID: enrollmentId,
		LessonID:     lessonId,
		Status:       models.ProgressStatusCompleted,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress updated successfully",
	})
}

func (h *ProgressHandler) getProgress(c *gin.Context) {
	enrollmentId, err := strconv.ParseInt(c.Param("enrollment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.UserID(c)

	progress, err := h.progressService.GetProgress(c, userID, enrollmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"progress": progress,
	})
}
