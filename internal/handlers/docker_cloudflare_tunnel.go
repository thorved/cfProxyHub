package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cfProxyHub/internal/services"
	"cfProxyHub/pkg/utils"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

// DockerCloudflareTunnelHandler handles HTTP requests related to Cloudflare Tunnels in Docker
type DockerCloudflareTunnelHandler struct {
	dockerService *services.DockerService
}

// NewDockerCloudflareTunnelHandler creates a new DockerCloudflareTunnelHandler instance
func NewDockerCloudflareTunnelHandler(dockerService *services.DockerService) *DockerCloudflareTunnelHandler {
	return &DockerCloudflareTunnelHandler{
		dockerService: dockerService,
	}
}

// NormalizeContainerResponse standardizes container response data
// It ensures that both ID and Id fields are present for frontend compatibility
func NormalizeContainerResponse(container types.Container) map[string]interface{} {
	// Create a map from the container to allow adding fields
	result := map[string]interface{}{
		"ID":       container.ID, // Add standard upper-case ID
		"Id":       container.ID, // Add standard lower-case Id
		"Names":    container.Names,
		"Image":    container.Image,
		"ImageID":  container.ImageID,
		"Command":  container.Command,
		"Created":  container.Created,
		"Ports":    container.Ports,
		"Labels":   container.Labels,
		"State":    container.State,
		"Status":   container.Status,
		"Networks": container.NetworkSettings.Networks,
	}

	return result
}

// ListTunnels returns all Cloudflare Tunnel containers
func (h *DockerCloudflareTunnelHandler) ListTunnels(c *gin.Context) {
	// Log the request for debugging
	c.Header("X-Debug", "ListTunnels endpoint called")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Add CORS header for debugging

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if Docker service is connected
	if h.dockerService == nil {
		utils.ErrorResponse(c, "Docker service not initialized", http.StatusInternalServerError)
		return
	}

	// First check Docker connectivity
	pingErr := h.dockerService.Ping(ctx)
	if pingErr != nil {
		c.Header("X-Debug-Docker", "Docker daemon connection failed")
		utils.ErrorResponse(c, "Failed to connect to Docker daemon: "+pingErr.Error(), http.StatusInternalServerError)
		return
	}

	// Try to find Cloudflare tunnel containers
	tunnels, err := h.dockerService.FindCloudflareTunnelContainers(ctx)
	if err != nil {
		// Log the error
		c.Error(err)
		c.Header("X-Debug-Error", err.Error())

		// Try listing all containers as a fallback
		allContainers, allErr := h.dockerService.ListContainers(ctx)
		if allErr != nil {
			c.Header("X-Debug-AllError", allErr.Error())
			utils.ErrorResponse(c, "Failed to list any containers: "+allErr.Error(), http.StatusInternalServerError)
			return
		}

		// Filter the containers for potential Cloudflare tunnels based on image name
		var possibleTunnels []types.Container
		for _, container := range allContainers {
			if strings.Contains(strings.ToLower(container.Image), "cloudflare") {
				possibleTunnels = append(possibleTunnels, container)
			}
		}

		// If we found any, return those
		if len(possibleTunnels) > 0 {
			c.Header("X-Debug-Info", "Returning filtered containers by image")
			utils.SuccessResponse(c, possibleTunnels)
			return
		}

		// Otherwise return all containers
		c.Header("X-Debug-Info", "Returning all containers as fallback")
		utils.SuccessResponse(c, allContainers)
		return
	}

	// If no tunnels found, return empty array instead of null
	if tunnels == nil || len(tunnels) == 0 {
		c.Header("X-Debug-Info", "No Cloudflare tunnel containers found")
		utils.SuccessResponse(c, []map[string]interface{}{})
		return
	}

	// Normalize container data for consistent ID field naming
	normalizedTunnels := make([]map[string]interface{}, len(tunnels))
	for i, tunnel := range tunnels {
		normalizedTunnels[i] = NormalizeContainerResponse(tunnel)
	}

	c.Header("X-Debug-Info", fmt.Sprintf("Found %d Cloudflare tunnel containers", len(tunnels)))
	utils.SuccessResponse(c, normalizedTunnels)
}

// CreateTunnel creates a new Cloudflare Tunnel container
func (h *DockerCloudflareTunnelHandler) CreateTunnel(c *gin.Context) {
	var params services.CloudflareTunnelParams
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if params.Token == "" {
		utils.ErrorResponse(c, "Tunnel token is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containerID, err := h.dockerService.CreateCloudflareTunnelContainer(ctx, params)
	if err != nil {
		utils.ErrorResponse(c, "Failed to create Cloudflare tunnel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Start the container
	if err := h.dockerService.StartContainer(ctx, containerID); err != nil {
		utils.ErrorResponse(c, "Tunnel created but failed to start: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"id":      containerID,
		"message": "Cloudflare tunnel created and started successfully",
	})
}

// DeleteTunnel removes a Cloudflare Tunnel container
func (h *DockerCloudflareTunnelHandler) DeleteTunnel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, "Container ID is required", http.StatusBadRequest)
		return
	}

	// Check if ID is "Unknown"
	if id == "Unknown" {
		utils.ErrorResponse(c, "Invalid container ID: Container not properly created or detected", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check Docker connectivity first
	pingErr := h.dockerService.Ping(ctx)
	if pingErr != nil {
		utils.ErrorResponse(c, "Docker daemon is not accessible: "+pingErr.Error(), http.StatusInternalServerError)
		return
	}

	// Try to inspect the container first to verify it exists
	_, err := h.dockerService.InspectContainer(ctx, id)
	if err != nil {
		// If container doesn't exist, that's actually OK for deletion
		// We'll just return success since the end result is what the user wanted
		utils.SuccessResponse(c, gin.H{
			"message": "Container doesn't exist or was already removed",
			"id":      id,
		})
		return
	}

	if err := h.dockerService.RemoveContainer(ctx, id); err != nil {
		utils.ErrorResponse(c, "Failed to remove Cloudflare tunnel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Cloudflare tunnel removed successfully",
		"id":      id,
	})
}

// StartTunnel starts a Cloudflare Tunnel container
func (h *DockerCloudflareTunnelHandler) StartTunnel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, "Container ID is required", http.StatusBadRequest)
		return
	}

	// Check if ID is "Unknown"
	if id == "Unknown" {
		utils.ErrorResponse(c, "Invalid container ID: Container not properly created or detected", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check Docker connectivity first
	pingErr := h.dockerService.Ping(ctx)
	if pingErr != nil {
		utils.ErrorResponse(c, "Docker daemon is not accessible: "+pingErr.Error(), http.StatusInternalServerError)
		return
	}

	// Try to inspect the container first to verify it exists
	_, err := h.dockerService.InspectContainer(ctx, id)
	if err != nil {
		utils.ErrorResponse(c, "Container not found or cannot be accessed: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := h.dockerService.StartContainer(ctx, id); err != nil {
		utils.ErrorResponse(c, "Failed to start Cloudflare tunnel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Cloudflare tunnel started successfully",
		"id":      id,
	})
}

// StopTunnel stops a Cloudflare Tunnel container
func (h *DockerCloudflareTunnelHandler) StopTunnel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, "Container ID is required", http.StatusBadRequest)
		return
	}

	// Check if ID is "Unknown"
	if id == "Unknown" {
		utils.ErrorResponse(c, "Invalid container ID: Container not properly created or detected", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check Docker connectivity first
	pingErr := h.dockerService.Ping(ctx)
	if pingErr != nil {
		utils.ErrorResponse(c, "Docker daemon is not accessible: "+pingErr.Error(), http.StatusInternalServerError)
		return
	}

	// Try to inspect the container first to verify it exists
	_, err := h.dockerService.InspectContainer(ctx, id)
	if err != nil {
		utils.ErrorResponse(c, "Container not found or cannot be accessed: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := h.dockerService.StopContainer(ctx, id); err != nil {
		utils.ErrorResponse(c, "Failed to stop Cloudflare tunnel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Cloudflare tunnel stopped successfully",
		"id":      id,
	})
}

// RestartTunnel restarts a Cloudflare Tunnel container
func (h *DockerCloudflareTunnelHandler) RestartTunnel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, "Container ID is required", http.StatusBadRequest)
		return
	}

	// Check if ID is "Unknown"
	if id == "Unknown" {
		utils.ErrorResponse(c, "Invalid container ID: Container not properly created or detected", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check Docker connectivity first
	pingErr := h.dockerService.Ping(ctx)
	if pingErr != nil {
		utils.ErrorResponse(c, "Docker daemon is not accessible: "+pingErr.Error(), http.StatusInternalServerError)
		return
	}

	// Try to inspect the container first to verify it exists
	_, err := h.dockerService.InspectContainer(ctx, id)
	if err != nil {
		utils.ErrorResponse(c, "Container not found or cannot be accessed: "+err.Error(), http.StatusNotFound)
		return
	}

	// Stop the container first
	if err := h.dockerService.StopContainer(ctx, id); err != nil {
		utils.ErrorResponse(c, "Failed to stop Cloudflare tunnel for restart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Start the container again
	if err := h.dockerService.StartContainer(ctx, id); err != nil {
		utils.ErrorResponse(c, "Failed to start Cloudflare tunnel after stop: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Cloudflare tunnel restarted successfully",
	})
}

// DockerDebugInfo returns diagnostic information about Docker
func (h *DockerCloudflareTunnelHandler) DockerDebugInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	debugInfo := map[string]interface{}{
		"docker_status": "unknown",
		"errors":        []string{},
		"containers":    []interface{}{},
		"timestamp":     time.Now().String(),
	}

	// Check Docker connectivity
	pingErr := h.dockerService.Ping(ctx)
	if pingErr != nil {
		debugInfo["docker_status"] = "disconnected"
		debugInfo["errors"] = append(debugInfo["errors"].([]string), "Docker connection error: "+pingErr.Error())
		utils.SuccessResponse(c, debugInfo)
		return
	}

	debugInfo["docker_status"] = "connected"

	// Get all containers
	containers, err := h.dockerService.ListContainers(ctx)
	if err != nil {
		debugInfo["errors"] = append(debugInfo["errors"].([]string), "Container listing error: "+err.Error())
	} else {
		// Add container info
		simplifiedContainers := []map[string]interface{}{}
		for _, container := range containers {
			simplifiedContainers = append(simplifiedContainers, map[string]interface{}{
				"id":     container.ID,
				"names":  container.Names,
				"image":  container.Image,
				"state":  container.State,
				"status": container.Status,
				"labels": container.Labels,
			})
		}
		debugInfo["containers"] = simplifiedContainers
		debugInfo["container_count"] = len(containers)
	}

	// Try to find Cloudflare tunnels
	tunnels, err := h.dockerService.FindCloudflareTunnelContainers(ctx)
	if err != nil {
		debugInfo["errors"] = append(debugInfo["errors"].([]string), "Tunnel detection error: "+err.Error())
	} else {
		debugInfo["tunnel_count"] = len(tunnels)
	}

	utils.SuccessResponse(c, debugInfo)
}
