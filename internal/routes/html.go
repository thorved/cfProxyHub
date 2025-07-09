package routes

import (
	"net/http"
	"path/filepath"

	"cfPorxyHub/internal/config"
	"cfPorxyHub/internal/handlers"
	"cfPorxyHub/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupHTMLRoutes configures HTML-related routes
func SetupHTMLRoutes(router *gin.Engine, cfg *config.Config) {
	// Dynamically load all HTML templates from the templates directory
	var templatePaths []string

	// Use multiple glob patterns to catch templates at different directory levels
	patterns := []string{
		"web/templates/*/*.html",   // 2 levels: auth/, dashboard/, layouts/
		"web/templates/*/*/*.html", // 3 levels: cloudflare/accounts/, cloudflare/tunnels/, cloudflare/zones/
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			panic("Failed to load template files with pattern " + pattern + ": " + err.Error())
		}
		templatePaths = append(templatePaths, matches...)
	}

	if len(templatePaths) == 0 {
		panic("No HTML templates found in web/templates/")
	}

	// Load all discovered templates
	router.LoadHTMLFiles(templatePaths...)

	// Serve static files
	router.Static("/assets", "./web/assets")

	// Initialize auth handler
	authHandler := handlers.NewLoginHandler(cfg)
	// Public routes (no authentication required)
	router.GET("/login", authHandler.LoginForm)
	router.GET("/logout", authHandler.Logout)

	// Protected routes (authentication required)
	router.GET("/", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "Dashboard.html", gin.H{})
	})

	// Cloudflare Account routes
	router.GET("/cloudflare/accounts", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareAccounts.html", gin.H{})
	})

	// Cloudflare Tunnel routes
	router.GET("/cloudflare/tunnels", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareAllTunnels.html", gin.H{})
	})

	router.GET("/cloudflare/tunnels/create", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "Cloudflare_CreateTunnel.html", gin.H{})
	})

	router.GET("/cloudflare/tunnels/hostnames", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "Cloudflare_TunnelPublicHostname.html", gin.H{})
	})

	// Cloudflare Zone routes
	router.GET("/cloudflare/zones", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareZones.html", gin.H{})
	})

	router.GET("/cloudflare/zones/details", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareZoneDetails.html", gin.H{})
	})

	// Legacy routes for backward compatibility (optional - can be removed later)
	router.GET("/CloudflareAccounts", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareAccounts.html", gin.H{})
	})
	router.GET("/CloudflareAllTunnels", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareAllTunnels.html", gin.H{})
	})
	router.GET("/Cloudflare_CreateTunnel", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "Cloudflare_CreateTunnel.html", gin.H{})
	})
	router.GET("/Cloudflare_TunnelPublicHostname", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "Cloudflare_TunnelPublicHostname.html", gin.H{})
	})
	router.GET("/CloudflareZones", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareZones.html", gin.H{})
	})
	router.GET("/CloudflareZoneDetails", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareZoneDetails.html", gin.H{})
	})

}
