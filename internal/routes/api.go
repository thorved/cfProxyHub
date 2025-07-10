package routes

import (
	"cfProxyHub/internal/config"
	"cfProxyHub/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes configures the API routes for the application (excluding auth routes)
func SetupAPIRoutes(router *gin.Engine, cfg *config.Config) {

	// Create a new group for API routes
	api := router.Group("/api")

	// Health check endpoint (public - no authentication required)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "cfProxyHub",
			"timestamp": gin.H{
				"unix": c.Request.Header.Get("X-Request-Time"),
			},
		})
	})

	// Apply authentication middleware to all API routes (except health and auth routes)
	api.Use(middleware.RequireAuth())

	// Status endpoint to check authentication
	api.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"message":       "User is authenticated",
		})
	})

}
