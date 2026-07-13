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
	Roles    string `json:"roles" binding:"required"`
}

// @Summary Đăng nhập tài khoản
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Tên đăng nhập và mật khẩu"
// @Success 200 {object} object{access_token=string} "Đăng nhập thành công"
// @Failure 400 {object} response.ErrorResponse "Dữ liệu không hợp lệ"
// @Failure 401 {object} response.ErrorResponse "Sai tên đăng nhập hoặc mật khẩu"
// @Failure 429 {object} response.ErrorResponse "Đăng nhập quá nhiều lần"
// @Router /auth/login [post]
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

// @Summary Đăng ký tài khoản mới
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Thông tin tài khoản"
// @Success 200 {object} object{message=string} "Đăng ký tài khoản thành công"
// @Failure 400 {object} response.ErrorResponse "Dữ liệu không hợp lệ"
// @Router /auth/register [post]
func (h *AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	_, err := h.userService.Register(c.Request.Context(), req.FullName, req.Username, req.Password, req.Roles)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}
