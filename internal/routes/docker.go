package routes

import (
	"cfProxyHub/internal/handlers"
	"cfProxyHub/internal/middleware"
	"cfProxyHub/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterDockerRoutes sets up Docker-related API endpoints
func RegisterDockerRoutes(api *gin.Engine) {
	dockerService, err := services.NewDockerService()
	if err != nil {
		return // Optionally log error
	}
	dockerHandler := handlers.NewDockerHandler(dockerService)

	docker := api.Group("/api/docker")
	docker.Use(middleware.RequireAuth())

	// Container management
	docker.GET("/containers", dockerHandler.ListContainers)
	docker.POST("/containers", dockerHandler.CreateContainer)
	docker.DELETE("/containers/:id", dockerHandler.RemoveContainer)
	docker.POST("/containers/:id/start", dockerHandler.StartContainer)
	docker.POST("/containers/:id/stop", dockerHandler.StopContainer)

	// Other Docker resources
	docker.GET("/images", dockerHandler.ListImages)
	docker.GET("/volumes", dockerHandler.ListVolumes)
	docker.GET("/networks", dockerHandler.ListNetworks)
}
