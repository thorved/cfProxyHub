package routes

import (
	"cfProxyHub/internal/config"
	"cfProxyHub/internal/handlers"
	"cfProxyHub/internal/middleware"
	"cfProxyHub/internal/services"
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
	zoneHandler := handlers.NewCloudflareZoneHandler(cfService)

	// Cloudflare API routes
	cloudflare := router.Group("/api/cloudflare")

	// Apply authentication middleware to all Cloudflare routes
	cloudflare.Use(middleware.RequireAuth())

	{
		// Account routes
		cloudflare.GET("/accounts", cfAccountHandler.GetAccounts)               // JSON API for direct API access
		cloudflare.GET("/accounts/:accountId", cfAccountHandler.GetAccountByID) // Get specific account by ID

		// Zone routes
		cloudflare.GET("/accounts/:accountId/zones", zoneHandler.GetZonesByAccountID)               // Get zones for specific account
		cloudflare.POST("/accounts/:accountId/zones", zoneHandler.CreateZone)                       // Create new zone
		cloudflare.GET("/zones/:zoneId", zoneHandler.GetZoneByID)                                   // Get specific zone by ID
		cloudflare.PUT("/zones/:zoneId", zoneHandler.UpdateZone)                                    // Update existing zone
		cloudflare.DELETE("/zones/:zoneId", zoneHandler.DeleteZone)                                 // Delete zone
		cloudflare.GET("/accounts/:accountId/zones/by-name/:domainName", zoneHandler.GetZoneByName) // Get zone by domain name
		cloudflare.GET("/accounts/:accountId/zones/dropdown", zoneHandler.GetZonesForDropdown)      // Get zones for dropdown usage

		// Tunnel routes
		cloudflare.GET("/accounts/:accountId/tunnels", tunnelHandler.GetTunnelsByAccountID)           // Get tunnels for specific account
		cloudflare.GET("/accounts/:accountId/tunnels/:tunnel_id", tunnelHandler.GetTunnelByID)        // Get specific tunnel by ID
		cloudflare.POST("/accounts/:accountId/tunnels", tunnelHandler.CreateTunnel)                   // Create new tunnel
		cloudflare.PUT("/accounts/:accountId/tunnels/:tunnel_id", tunnelHandler.UpdateTunnel)         // Update existing tunnel
		cloudflare.DELETE("/accounts/:accountId/tunnels/:tunnel_id", tunnelHandler.DeleteTunnel)      // Delete tunnel
		cloudflare.GET("/accounts/:accountId/tunnels/:tunnel_id/token", tunnelHandler.GetTunnelToken) // Get tunnel token

		// Public hostname routes
		cloudflare.GET("/accounts/:accountId/tunnels/:tunnel_id/hostnames", tunnelHandler.GetPublicHostnamesByTunnelID)      // Get public hostnames for tunnel
		cloudflare.POST("/accounts/:accountId/tunnels/:tunnel_id/hostnames", tunnelHandler.CreatePublicHostname)             // Create public hostname
		cloudflare.PUT("/accounts/:accountId/tunnels/:tunnel_id/hostnames/:hostname", tunnelHandler.UpdatePublicHostname)    // Update public hostname
		cloudflare.DELETE("/accounts/:accountId/tunnels/:tunnel_id/hostnames/:hostname", tunnelHandler.DeletePublicHostname) // Delete public hostname
	}
}
