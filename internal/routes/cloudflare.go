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
	cfAccountHandler := handlers.NewCloudflareAccountHandler(cfService)
	tunnelHandler := handlers.NewCloudflareTunnelHandler(cfService)

	// Cloudflare API routes
	cloudflare := router.Group("/api/cloudflare")

	// Apply authentication middleware to all Cloudflare routes
	cloudflare.Use(middleware.RequireAuth())

	{
		// Account routes
		cloudflare.GET("/accounts", cfAccountHandler.GetAccounts)        // JSON API for direct API access
		cloudflare.GET("/accounts/:id", cfAccountHandler.GetAccountByID) // Get specific account by ID

		// Tunnel routes
		cloudflare.GET("/accounts/:id/tunnels", tunnelHandler.GetTunnelsByAccountID)           // Get tunnels for specific account
		cloudflare.GET("/accounts/:id/tunnels/:tunnel_id", tunnelHandler.GetTunnelByID)        // Get specific tunnel by ID
		cloudflare.POST("/accounts/:id/tunnels", tunnelHandler.CreateTunnel)                   // Create new tunnel
		cloudflare.PUT("/accounts/:id/tunnels/:tunnel_id", tunnelHandler.UpdateTunnel)         // Update existing tunnel
		cloudflare.DELETE("/accounts/:id/tunnels/:tunnel_id", tunnelHandler.DeleteTunnel)      // Delete tunnel
		cloudflare.GET("/accounts/:id/tunnels/:tunnel_id/token", tunnelHandler.GetTunnelToken) // Get tunnel token
	}
}
