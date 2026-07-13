package response

import (
	"errors"
	"learning-platform/internal/models"
	"net/http"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Error(c *gin.Context, err error) {
	var appErr *models.AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus >= 500 {
			logger.Error("Internal server error:",
				"code", appErr.Code,
				"error", appErr.Error,
				"path", c.Request.URL.Path,
			)
		}

		c.JSON(appErr.HTTPStatus, ErrorResponse{
			Error: ErrorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
		return
	}

	logger.Error("Error:", err)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
		},
	})
}

func AbortWithError(c *gin.Context, err error) {
	var appErr *models.AppError
	if errors.As(err, &appErr) {
		c.AbortWithStatusJSON(appErr.HTTPStatus, ErrorResponse{
			Error: ErrorDetail{Code: appErr.Code, Message: appErr.Message},
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal error"},
	})
}
