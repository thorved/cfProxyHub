package routes

import (
	"cfProxyHub/internal/config"
	"cfProxyHub/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures the authentication API routes
func SetupAuthRoutes(router *gin.Engine, cfg *config.Config) {
	// Create auth route group under /api/auth
	auth := router.Group("/api/auth")

	// Initialize auth handler
	authHandler := handlers.NewLoginHandler(cfg)

	// Authentication endpoints (public - no auth required)
	auth.POST("/login", authHandler.LoginAPI)
	auth.POST("/logout", authHandler.LogoutAPI)
}
