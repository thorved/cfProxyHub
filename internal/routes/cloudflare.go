package routes

import (
	"cfPorxyHub/internal/config"
	"cfPorxyHub/internal/handlers"
	"cfPorxyHub/internal/middleware"
	"cfPorxyHub/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes configures the API routes for the application
func SetupCloudflareRoutes(router *gin.Engine) {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Cloudflare service
	cfService, err := services.NewCloudflareService(cfg.CloudflareAPIToken, cfg.CloudflareAPIKey, cfg.CloudflareEmail)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudflare service: %v", err)
	}

	// Initialize handlers
	cfHandler := handlers.NewCloudflareHandler(cfService)

	// Cloudflare API routes
	cloudflare := router.Group("/api/cloudflare")

	// Apply authentication middleware to all Cloudflare routes
	cloudflare.Use(middleware.RequireAuth())

	{
		cloudflare.GET("/accounts/json", cfHandler.GetAccounts) // JSON API for direct API access
		cloudflare.GET("/accounts/:id", cfHandler.GetAccountByID)
	}
}
