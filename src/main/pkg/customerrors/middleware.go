package customerrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler is a middleware that handles API errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case Error:
				c.JSON(e.StatusCode(), ErrorResponse{
					Code:    e.ErrorCode(),
					Message: e.Error(),
				})
			case *ValidationErrorResponse:
				c.JSON(e.StatusCode(), e)
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "An unexpected error occurred",
				})
			}
			return
		}
	}
}
