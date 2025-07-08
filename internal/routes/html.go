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
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/accounts", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "accounts.html", gin.H{})
	})
}
