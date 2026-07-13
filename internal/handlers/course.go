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
		courseGroup.GET("/:course_id", authMW, h.GetCourseByID)
		courseGroup.PATCH("/:course_id", authMW, h.update)
		courseGroup.DELETE("/:course_id", authMW, h.delete)
	}
}

type courseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TotalSeats  int    `json:"total_seats"`
}

type courseRequestFilter struct {
	PageID  int    `form:"page_id"`
	PerPage int    `form:"per_page"`
	Status  string `form:"status"`
	Keyword string `form:"keyword"`
}

type updateCourseRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	TotalSeats  *int    `json:"total_seats"`
	Status      *string `json:"status"`
}

func (h *CourseHandler) create(c *gin.Context) {
	var req courseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	err := h.courseService.CreateCourse(c.Request.Context(), &models.Course{
		InstructorID: middleware.UserID(c),
		Title:        req.Title,
		Description:  req.Description,
		Status:       models.CourseStatusPublished,
		TotalSeats:   req.TotalSeats,
	})

	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Course created",
	})
}

func (h *CourseHandler) list(c *gin.Context) {
	var req courseRequestFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err)
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
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	course, err := h.courseService.GetCourseByID(c.Request.Context(), courseID)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"course": course,
	})
}

func (h *CourseHandler) update(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req updateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
	}

	input := &models.UpdateCourse{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		TotalSeats:  req.TotalSeats,
	}

	err = h.courseService.UpdateCourse(c.Request.Context(), middleware.UserID(c), courseID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course updated",
	})
}

func (h *CourseHandler) delete(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.courseService.DeleteCourse(c.Request.Context(), middleware.UserID(c), courseID)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course deleted",
	})
}
