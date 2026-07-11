package handlers

import (
	"learning-platform/internal/models"
	"learning-platform/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	lessonService services.LessonService
}

func NewLessonHandler(lessonService services.LessonService) *LessonHandler {
	return &LessonHandler{
		lessonService: lessonService,
	}
}

func (h *LessonHandler) RegisterRoutes(r *gin.RouterGroup, authMW gin.HandlerFunc) {
	lessonGroup := r.Group("/lessons")
	{
		lessonGroup.GET("", authMW, h.list)
		lessonGroup.POST("", authMW, h.create)
	}
}

type lessonRequestFilter struct {
	PageID   int   `form:"page_id"`
	PerPage  int   `form:"per_page"`
	CourseID int64 `form:"course_id" binding:"required"`
}

type lessonRequest struct {
	CourseID   int64  `json:"course_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	OrderIndex int64  `json:"order_index"`
}

func (h *LessonHandler) list(c *gin.Context) {
	var req lessonRequestFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := services.LessonFilter{
		PageID:   req.PageID,
		PerPage:  req.PerPage,
		CourseID: req.CourseID,
	}
	courses, err := h.lessonService.List(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (h *LessonHandler) create(c *gin.Context) {
	var req lessonRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.lessonService.CreateLesson(c.Request.Context(), &models.Lesson{
		CourseID:   req.CourseID,
		Title:      req.Title,
		Content:    req.Content,
		OrderIndex: req.OrderIndex,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lesson created",
	})
}
