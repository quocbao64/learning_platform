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
