package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}

func ErrorResponse(c *gin.Context, message string, code int) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}
