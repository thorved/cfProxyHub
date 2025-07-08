package routes

import (
	"cfPorxyHub/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes configures the API routes for the application
func SetupAPIRoutes(router *gin.Engine) {
	// Create a new group for API routes
	api := router.Group("/api")

	// Apply authentication middleware to all API routes
	api.Use(middleware.RequireAuth())

	// Define a simple GET endpoint (now protected)
	api.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
}
