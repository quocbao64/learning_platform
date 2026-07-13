package handlers

import (
	"learning-platform/internal/handlers/response"
	"learning-platform/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthHandler(authService services.AuthService, userService services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup, ratelimit gin.HandlerFunc) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", ratelimit, h.login)
		authGroup.POST("/register", h.register)
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

func (h *AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	_, err := h.userService.Register(c.Request.Context(), req.FullName, req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}
