package handlers

import (
	"learning-platform/internal/middleware"
	"learning-platform/internal/models"
	"learning-platform/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService services.CourseService
}

func NewCourseHandler(courseService services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

func (h *CourseHandler) RegisterRoute(r *gin.RouterGroup, authMW gin.HandlerFunc) {
	courseGroup := r.Group("/courses")
	{
		courseGroup.POST("", authMW, h.create)
		courseGroup.GET("", authMW, h.list)
		courseGroup.GET("/:id", authMW, h.GetCourseByID)
	}
}

type courseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type courseRequestFilter struct {
	PageID  int    `form:"page_id,default=1"`
	PerPage int    `form:"per_page,default=10"`
	Status  string `form:"status"`
	Keyword string `form:"keyword"`
}

func (h *CourseHandler) create(c *gin.Context) {
	var req courseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.courseService.CreateCourse(c.Request.Context(), &models.Course{
		InstructorID: middleware.UserID(c),
		Title:        req.Title,
		Description:  req.Description,
		Status:       models.CourseStatusDraft,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Course created",
	})
}

func (h *CourseHandler) list(c *gin.Context) {
	var req courseRequestFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := &services.CourseFilter{
		PageID:  req.PageID,
		PerPage: req.PerPage,
		Status:  req.Status,
		Keyword: req.Keyword,
	}
	courses, err := h.courseService.ListCourses(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseService.GetCourseByID(c.Request.Context(), int64(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"course": course,
	})
}
