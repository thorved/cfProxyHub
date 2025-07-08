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
func SetupCloudflareRoutes(router *gin.Engine, cfg *config.Config) {

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
		cloudflare.GET("/accounts", cfHandler.GetAccounts)                                 // JSON API for direct API access
		cloudflare.GET("/accounts/:id", cfHandler.GetAccountByID)                          // Get specific account by ID
		cloudflare.GET("/accounts/:id/tunnels", cfHandler.GetTunnelsByAccountID)           // Get tunnels for specific account
		cloudflare.GET("/accounts/:id/tunnels/:tunnel_id", cfHandler.GetTunnelByID)        // Get specific tunnel by ID
		cloudflare.POST("/accounts/:id/tunnels", cfHandler.CreateTunnel)                   // Create new tunnel
		cloudflare.PUT("/accounts/:id/tunnels/:tunnel_id", cfHandler.UpdateTunnel)         // Update existing tunnel
		cloudflare.DELETE("/accounts/:id/tunnels/:tunnel_id", cfHandler.DeleteTunnel)      // Delete tunnel
		cloudflare.GET("/accounts/:id/tunnels/:tunnel_id/token", cfHandler.GetTunnelToken) // Get tunnel token
	}
}
