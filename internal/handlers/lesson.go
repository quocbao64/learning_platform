package handlers

import (
	"learning-platform/internal/models"
	"learning-platform/internal/services"
	"net/http"
	"strconv"
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
	lessonGroup := r.Group("/courses/:course_id/lessons")
	{
		lessonGroup.GET("", authMW, h.list)
		lessonGroup.POST("", authMW, h.create)
	}
}

type lessonRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	OrderIndex int64  `json:"order_index" binding:"required"`
}

type listLessonsRequest struct {
	PageID  int `form:"page_id"`
	PerPage int `form:"per_page"`
}

func (h *LessonHandler) list(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req listLessonsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := services.LessonFilter{
		CourseID: courseID,
		PageID:   req.PageID,
		PerPage:  req.PerPage,
	}
	courses, err := h.lessonService.List(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": courses,
	})
}

func (h *LessonHandler) create(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req lessonRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.lessonService.CreateLesson(c.Request.Context(), &models.Lesson{
		CourseID:   courseID,
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
