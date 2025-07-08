package routes

import (
	"net/http"

	"cfPorxyHub/internal/config"
	"cfPorxyHub/internal/handlers"
	"cfPorxyHub/internal/middleware"
	"cfPorxyHub/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupHTMLRoutes configures HTML-related routes
func SetupHTMLRoutes(router *gin.Engine, cfg *config.Config) {
	// Load HTML templates including partials subfolder
	router.LoadHTMLGlob("web/templates/*")

	// Serve static files
	router.Static("/assets", "./web/assets")

	// Initialize auth handler
	authHandler := handlers.NewLoginHandler(cfg)
	// Public routes (no authentication required)
	router.GET("/login", authHandler.LoginForm)
	router.POST("/login", authHandler.Login)
	router.GET("/logout", authHandler.Logout)

	// HTMX component routes
	components := router.Group("/components")
	components.Use(middleware.AuthMiddleware())
	// Component routes for HTMX - updated paths
	components.GET("/header", func(c *gin.Context) {
		c.HTML(http.StatusOK, "header.html", gin.H{})
	})
	components.GET("/sidebar", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sidebar.html", gin.H{})
	})
	components.GET("/footer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "footer.html", gin.H{})
	})

	// Protected routes (authentication required)
	router.GET("/", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "Dashboard.html", gin.H{})
	})

	router.GET("/CloudflareAccounts", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "CloudflareAccounts.html", gin.H{})
	})

	// HTMX API endpoints for dynamic content loading
	htmxAPI := router.Group("/htmx")
	htmxAPI.Use(middleware.AuthMiddleware())

	// Initialize Cloudflare service and handler for HTMX endpoints
	cfService, err := services.NewCloudflareService(cfg.CloudflareAPIToken, cfg.CloudflareAPIKey, cfg.CloudflareEmail)
	if err != nil {
		// Handle error appropriately in production
		panic("Failed to initialize Cloudflare service: " + err.Error())
	}
	cfHandler := handlers.NewCloudflareHandler(cfService)

	// HTMX endpoints
	htmxAPI.GET("/accounts", cfHandler.GetAccountsHTML)
}
