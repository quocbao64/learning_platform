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
	instructorRoles := middleware.RequireRole(models.UserRoleInstructor, models.UserRoleAdmin)
	courseGroup := r.Group("/courses")
	{
		courseGroup.POST("", authMW, instructorRoles, h.create)
		courseGroup.GET("", authMW, h.list)
		courseGroup.GET("/:course_id", authMW, h.getCourseByID)
		courseGroup.PATCH("/:course_id", authMW, instructorRoles, h.update)
		courseGroup.DELETE("/:course_id", authMW, instructorRoles, h.delete)
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

// @Summary Tạo khóa học mới
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courseRequest true "Thông tin khóa học"
// @Success 201 {object} object{message=string} "Tạo khóa học thành công"
// @Failure 400 {object} response.ErrorResponse "Dữ liệu đầu vào không hợp lệ"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Failure 403 {object} response.ErrorResponse "Không có quyền tạo khóa học"
// @Router /courses [post]
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

// @Summary Danh sách khóa học
// @Tags Courses
// @Produce json
// @Security BearerAuth
// @Param page_id query int false "Số trang"
// @Param per_page query int false "Số lượng mỗi trang"
// @Param status query string false "Lọc theo trạng thái (draft, published, archived)"
// @Param keyword query string false "Lọc theo từ khóa"
// @Success 200 {object} object{courses=[]models.Course} "Danh sách khóa học"
// @Failure 400 {object} response.ErrorResponse "Params không hợp lệ"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /courses [get]
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

// @Summary Chi tiết khóa học
// @Tags Courses
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "ID khóa học"
// @Success 200 {object} object{course=models.Course}
// @Failure 400 {object} response.ErrorResponse "Không tìm thấy khóa học"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Router /courses/{course_id} [get]
func (h *CourseHandler) getCourseByID(c *gin.Context) {
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

// @Summary Cập nhật khóa học
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "ID khóa học"
// @Param request body updateCourseRequest true "Thông tin cập nhật khóa học"
// @Success 200 {object} object{message=string} "Cập nhật khóa học thành công"
// @Failure 400 {object} response.ErrorResponse "Dữ liệu đầu vào không hợp lệ hoặc không tìm thấy khóa học"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Failure 403 {object} response.ErrorResponse "Không có quyền cập nhật khóa học"
// @Router /courses/{course_id} [patch]
func (h *CourseHandler) update(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req updateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
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

// @Summary Xóa khóa học
// @Tags Courses
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "ID khóa học"
// @Success 200 {object} object{message=string} "Xóa khóa học thành công"
// @Failure 400 {object} response.ErrorResponse "Không tìm thấy khóa học hoặc khóa học đã có người đăng ký"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Failure 403 {object} response.ErrorResponse "Không có quyền xóa khóa học"
// @Router /courses/{course_id} [delete]
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
