package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes configures the API routes for the application
func SetupAPIRoutes(router *gin.Engine) {
	// Create a new group for API routes
	api := router.Group("/api")

	// Define a simple GET endpoint
	api.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

}
