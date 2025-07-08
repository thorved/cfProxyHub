package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(router *gin.Engine) {
	// Load HTML templates
	router.LoadHTMLGlob("web/templates/*")

	// Serve static files
	router.Static("/assets", "./web/assets")
	// Setup different route groups
	SetupAPIRoutes(router)  // JSON API endpoints
	SetupHTMLRoutes(router) // HTML pages
}
