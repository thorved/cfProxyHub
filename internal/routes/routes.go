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
	SetupAuthRoutes(router, cfg)       // Authentication API endpoints (/api/auth/*)
	SetupAPIRoutes(router, cfg)        // Protected JSON API endpoints (/api/*)
	SetupHTMLRoutes(router, cfg)       // HTML pages
	SetupCloudflareRoutes(router, cfg) // Cloudflare-specific API endpoints
}
