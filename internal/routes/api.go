package routes

import (
	"cfPorxyHub/internal/config"
	"cfPorxyHub/internal/handlers"
	"cfPorxyHub/internal/middleware"
	"cfPorxyHub/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes configures the API routes for the application (excluding auth routes)
func SetupAPIRoutes(router *gin.Engine, cfg *config.Config) {

	// Create services
	accountService := services.NewAccountService()

	// Create handlers
	accountHandler := handlers.NewAccountHandler(accountService)

	// Create a new group for API routes
	api := router.Group("/api")

	// Apply authentication middleware to all API routes (except auth routes which are handled separately)
	api.Use(middleware.RequireAuth())

	// Status endpoint to check authentication
	api.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"message":       "User is authenticated",
		})
	})

	// Current account endpoints
	api.GET("/current-account", accountHandler.GetCurrentAccount)
	api.POST("/current-account", accountHandler.SetCurrentAccount)
	api.DELETE("/current-account", accountHandler.ClearCurrentAccount)
}
