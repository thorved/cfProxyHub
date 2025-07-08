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
	// Load HTML templates - Go templates will handle includes via {{template}} syntax
	router.LoadHTMLGlob("web/templates/*.html")

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

}
