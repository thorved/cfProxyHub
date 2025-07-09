package routes

import (
	"net/http"

	"cfPorxyHub/internal/config"
	"cfPorxyHub/internal/handlers"
	"cfPorxyHub/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupHTMLRoutes configures HTML-related routes
func SetupHTMLRoutes(router *gin.Engine, cfg *config.Config) {
	// Load HTML templates from all subdirectories using multiple glob patterns
	router.LoadHTMLGlob("web/templates/*/*.html")

	// Also load the root level templates if any exist
	// router.LoadHTMLFiles() can be used for specific files if needed

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

}
