package routes

import (
	"cfPorxyHub/internal/config"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(router *gin.Engine) {
	// Load config
	cfg := config.LoadConfig()

	// Setup different route groups
	SetupAPIRoutes(router)        // JSON API endpoints
	SetupHTMLRoutes(router, cfg)  // HTML pages
	SetupCloudflareRoutes(router) // Cloudflare-specific API endpoints
}
