package handlers

import (
	"learning-platform/internal/handlers/response"
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

// @Summary Danh sách bài học của khóa học
// @Tags Lessons
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "ID khóa học"
// @Param page_id query int false "Số trang"
// @Param per_page query int false "Số lượng mỗi trang"
// @Success 200 {object} object{data=[]models.Lesson} "Danh sách bài học"
// @Failure 400 {object} response.ErrorResponse "ID khóa học không hợp lệ"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /courses/{course_id}/lessons [get]
func (h *LessonHandler) list(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req listLessonsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err)
		return
	}

	filter := services.LessonFilter{
		CourseID: courseID,
		PageID:   req.PageID,
		PerPage:  req.PerPage,
	}
	courses, err := h.lessonService.List(c.Request.Context(), &filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": courses,
	})
}

// @Summary Tạo bài học mới
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "ID khóa học"
// @Param request body lessonRequest true "Thông tin bài học"
// @Success 201 {object} object{message=string} "Tạo bài học thành công"
// @Failure 400 {object} response.ErrorResponse "Dữ liệu đầu vào không hợp lệ"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /courses/{course_id}/lessons [post]
func (h *LessonHandler) create(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req lessonRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, err)
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
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lesson created",
	})
}
