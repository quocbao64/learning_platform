package handlers

import (
	"learning-platform/internal/handlers/response"
	"learning-platform/internal/middleware"
	"learning-platform/internal/models"
	"learning-platform/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService,
	}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	userGroup := rg.Group("users")
	{
		userGroup.GET("/me", authMiddleware, h.getCurrentUser)
	}
}

type userResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Roles    string `json:"roles"`
}

func toUserResponse(user *models.User) *userResponse {
	return &userResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Roles:    user.Roles,
	}
}

// @Summary Lấy thông tin tài khoản đang đăng nhập
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{user=userResponse} "Thông tin tài khoản"
// @Failure 401 {object} response.ErrorResponse "Chưa đăng nhập hoặc token không hợp lệ"
// @Failure 404 {object} response.ErrorResponse "Không tìm thấy người dùng"
// @Router /users/me [get]
func (h *UserHandler) getCurrentUser(c *gin.Context) {
	user, err := h.userService.GetByID(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": toUserResponse(user),
	})
}
