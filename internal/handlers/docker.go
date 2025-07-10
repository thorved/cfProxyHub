package handlers

import (
	"context"
	"net/http"
	"time"

	"cfProxyHub/internal/services"
	"cfProxyHub/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DockerHandler struct {
	service *services.DockerService
}

func NewDockerHandler(service *services.DockerService) *DockerHandler {
	return &DockerHandler{service: service}
}

func (h *DockerHandler) ListContainers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	containers, err := h.service.ListContainers(ctx)
	if err != nil {
		utils.ErrorResponse(c, "Failed to list containers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, containers)
}

func (h *DockerHandler) ListImages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	images, err := h.service.ListImages(ctx)
	if err != nil {
		utils.ErrorResponse(c, "Failed to list images: "+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, images)
}

func (h *DockerHandler) ListVolumes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	volumes, err := h.service.ListVolumes(ctx)
	if err != nil {
		utils.ErrorResponse(c, "Failed to list volumes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, volumes)
}

func (h *DockerHandler) ListNetworks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	networks, err := h.service.ListNetworks(ctx)
	if err != nil {
		utils.ErrorResponse(c, "Failed to list networks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, networks)
}

// CreateContainer handles container creation requests
func (h *DockerHandler) CreateContainer(c *gin.Context) {
	var params services.CreateContainerParams
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id, err := h.service.CreateContainer(ctx, params)
	if err != nil {
		utils.ErrorResponse(c, "Failed to create container: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"id":      id,
		"message": "Container created successfully",
	})
}

// RemoveContainer handles container removal requests
func (h *DockerHandler) RemoveContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, "Container ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.service.RemoveContainer(ctx, id); err != nil {
		utils.ErrorResponse(c, "Failed to remove container: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Container removed successfully",
	})
}

// StartContainer handles container start requests
func (h *DockerHandler) StartContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, "Container ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.service.StartContainer(ctx, id); err != nil {
		utils.ErrorResponse(c, "Failed to start container: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Container started successfully",
	})
}

// StopContainer handles container stop requests
func (h *DockerHandler) StopContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, "Container ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.service.StopContainer(ctx, id); err != nil {
		utils.ErrorResponse(c, "Failed to stop container: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Container stopped successfully",
	})
}
