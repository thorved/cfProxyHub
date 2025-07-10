package routes

import (
	"context"
	"net/http"
	"time"

	"cfProxyHub/internal/handlers"
	//"cfProxyHub/internal/middleware"
	"cfProxyHub/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterDockerCloudflareTunnelRoutes sets up Docker-based Cloudflare Tunnel API endpoints
func RegisterDockerCloudflareTunnelRoutes(api *gin.Engine) {
	dockerService, err := services.NewDockerService()
	if err != nil {
		return // Optionally log error
	}
	dockerCFTunnelHandler := handlers.NewDockerCloudflareTunnelHandler(dockerService)

	// Add debug middleware to log all requests
	api.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/api/docker/cloudflare/tunnels" {
			c.Header("X-Debug-Route", "Docker Cloudflare Tunnel route accessed")
			c.Next()
		} else {
			c.Next()
		}
	})

	// Add Docker debug endpoints to check connection and get detailed diagnostics
	api.GET("/api/docker/debug", func(c *gin.Context) {
		if dockerService == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Docker service is nil",
			})
			return
		}

		// Try to ping the Docker daemon
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pingErr := dockerService.Ping(ctx)
		if pingErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to connect to Docker: " + pingErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Successfully connected to Docker",
		})
	})

	// Add a more detailed diagnostics endpoint
	api.GET("/api/docker/diagnostics", dockerCFTunnelHandler.DockerDebugInfo)

	// Create a group for Docker-based Cloudflare Tunnel endpoints
	dockerTunnels := api.Group("/api/docker/cloudflare/tunnels")
	//dockerTunnels.Use(middleware.RequireAuth())

	// Tunnel management endpoints
	dockerTunnels.GET("", dockerCFTunnelHandler.ListTunnels)
	dockerTunnels.POST("", dockerCFTunnelHandler.CreateTunnel)
	dockerTunnels.DELETE("/:id", dockerCFTunnelHandler.DeleteTunnel)
	dockerTunnels.POST("/:id/start", dockerCFTunnelHandler.StartTunnel)
	dockerTunnels.POST("/:id/stop", dockerCFTunnelHandler.StopTunnel)
	dockerTunnels.POST("/:id/restart", dockerCFTunnelHandler.RestartTunnel)
}
