package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cfPorxyHub/internal/models"
	"cfPorxyHub/internal/services"
	"cfPorxyHub/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CloudflareTunnelHandler struct {
	cfService *services.CloudflareService
}

// NewCloudflareTunnelHandler creates a new Cloudflare Tunnel handler
func NewCloudflareTunnelHandler(cfService *services.CloudflareService) *CloudflareTunnelHandler {
	return &CloudflareTunnelHandler{
		cfService: cfService,
	}
}

// GetTunnelsByAccountID handles the GET /api/cloudflare/accounts/:accountId/tunnels endpoint
func (h *CloudflareTunnelHandler) GetTunnelsByAccountID(c *gin.Context) {
	accountID := c.Param("accountId")
	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tunnels, err := h.cfService.GetCloudflareTunnels(ctx, accountID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch tunnels: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Tunnels retrieved successfully",
		"account_id": accountID,
		"tunnels":    tunnels,
	})
}

// GetTunnelByID handles the GET /api/cloudflare/accounts/:accountId/tunnels/:tunnel_id endpoint
func (h *CloudflareTunnelHandler) GetTunnelByID(c *gin.Context) {
	accountID := c.Param("accountId")
	tunnelID := c.Param("tunnel_id")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tunnel, err := h.cfService.GetCloudflareTunnelByID(ctx, accountID, tunnelID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch tunnel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Tunnel retrieved successfully",
		"account_id": accountID,
		"tunnel":     tunnel,
	})
}

// CreateTunnel handles the POST /api/cloudflare/accounts/:accountId/tunnels endpoint
func (h *CloudflareTunnelHandler) CreateTunnel(c *gin.Context) {
	accountID := c.Param("accountId")
	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}

	// Create a simple struct to receive the JSON data
	var requestData struct {
		Name      string `json:"name"`
		ConfigSrc string `json:"config_src,omitempty"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.ErrorResponse(c, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if requestData.Name == "" {
		utils.ErrorResponse(c, "Tunnel name is required", http.StatusBadRequest)
		return
	}

	// Create the request using the helper function
	request := models.NewTunnelCreateRequest(requestData.Name, requestData.ConfigSrc)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tunnel, err := h.cfService.CreateCloudflareTunnel(ctx, accountID, request)
	if err != nil {
		utils.ErrorResponse(c, "Failed to create tunnel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Tunnel created successfully",
		"account_id": accountID,
		"tunnel":     tunnel,
	})
}

// UpdateTunnel handles the PUT /api/cloudflare/accounts/:accountId/tunnels/:tunnel_id endpoint
func (h *CloudflareTunnelHandler) UpdateTunnel(c *gin.Context) {
	accountID := c.Param("accountId")
	tunnelID := c.Param("tunnel_id")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}

	// Create a simple struct to receive the JSON data
	var requestData struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.ErrorResponse(c, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create the request using the helper function
	request := models.NewTunnelUpdateRequest(requestData.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tunnel, err := h.cfService.UpdateCloudflareTunnel(ctx, accountID, tunnelID, request)
	if err != nil {
		utils.ErrorResponse(c, "Failed to update tunnel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Tunnel updated successfully",
		"account_id": accountID,
		"tunnel":     tunnel,
	})
}

// DeleteTunnel handles the DELETE /api/cloudflare/accounts/:accountId/tunnels/:tunnel_id endpoint
func (h *CloudflareTunnelHandler) DeleteTunnel(c *gin.Context) {
	accountID := c.Param("accountId")
	tunnelID := c.Param("tunnel_id")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := h.cfService.DeleteCloudflareTunnel(ctx, accountID, tunnelID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to delete tunnel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Tunnel deleted successfully",
		"account_id": accountID,
		"tunnel_id":  tunnelID,
	})
}

// GetTunnelToken handles the GET /api/cloudflare/accounts/:accountId/tunnels/:tunnel_id/token endpoint
func (h *CloudflareTunnelHandler) GetTunnelToken(c *gin.Context) {
	accountID := c.Param("accountId")
	tunnelID := c.Param("tunnel_id")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := h.cfService.GetCloudflareTunnelToken(ctx, accountID, tunnelID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to get tunnel token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Tunnel token retrieved successfully",
		"account_id": accountID,
		"tunnel_id":  tunnelID,
		"token":      token,
	})
}

// GetPublicHostnamesByTunnelID handles the GET /api/cloudflare/accounts/:accountId/tunnels/:tunnel_id/hostnames endpoint
func (h *CloudflareTunnelHandler) GetPublicHostnamesByTunnelID(c *gin.Context) {
	accountID := c.Param("accountId")
	tunnelID := c.Param("tunnel_id")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hostnames, err := h.cfService.GetCloudflareTunnelPublicHostnames(ctx, accountID, tunnelID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch public hostnames: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Public hostnames retrieved successfully",
		"account_id": accountID,
		"tunnel_id":  tunnelID,
		"hostnames":  hostnames,
		"total":      len(hostnames),
	})
}

// CreatePublicHostname handles the POST /api/cloudflare/accounts/:accountId/tunnels/:tunnel_id/hostnames endpoint
func (h *CloudflareTunnelHandler) CreatePublicHostname(c *gin.Context) {
	accountID := c.Param("accountId")
	tunnelID := c.Param("tunnel_id")

	fmt.Printf("CreatePublicHostname called with accountID=%s, tunnelID=%s\n", accountID, tunnelID)

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}

	// Create a simple struct to receive the JSON data
	var requestData struct {
		Hostname string `json:"hostname"`
		Service  string `json:"service"`
		Path     string `json:"path,omitempty"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.ErrorResponse(c, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if requestData.Hostname == "" {
		utils.ErrorResponse(c, "Hostname is required", http.StatusBadRequest)
		return
	}
	if requestData.Service == "" {
		utils.ErrorResponse(c, "Service is required", http.StatusBadRequest)
		return
	}

	// Create the hostname parameter using the helper function
	hostnameParam := models.NewPublicHostnameIngressParam(requestData.Hostname, requestData.Service, requestData.Path)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use the new function that creates both tunnel config and DNS record
	result, err := h.cfService.CreateCloudflareTunnelPublicHostnameWithDNS(ctx, accountID, tunnelID, requestData.Hostname, hostnameParam)
	if err != nil {
		// Log the full error for debugging
		fmt.Printf("Error creating public hostname: %v\n", err)
		fmt.Printf("Request data: hostname=%s, service=%s, path=%s\n", requestData.Hostname, requestData.Service, requestData.Path)
		fmt.Printf("Account ID: %s, Tunnel ID: %s\n", accountID, tunnelID)

		// Provide more specific error messages
		errorMessage := "Failed to create public hostname"
		if strings.Contains(err.Error(), "already exists") {
			errorMessage = "A hostname with this name already exists"
		} else if strings.Contains(err.Error(), "invalid") {
			errorMessage = "Invalid hostname configuration"
		} else if strings.Contains(err.Error(), "unauthorized") {
			errorMessage = "Unauthorized access to this domain"
		} else if strings.Contains(err.Error(), "domain validation") {
			errorMessage = "Domain validation failed: " + err.Error()
		}

		utils.ErrorResponse(c, errorMessage+": "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Public hostname created successfully",
		"account_id": accountID,
		"tunnel_id":  tunnelID,
		"hostname":   requestData.Hostname,
		"config":     result,
	})
}

// UpdatePublicHostname handles the PUT /api/cloudflare/accounts/:accountId/tunnels/:tunnel_id/hostnames/:hostname endpoint
func (h *CloudflareTunnelHandler) UpdatePublicHostname(c *gin.Context) {
	accountID := c.Param("accountId")
	tunnelID := c.Param("tunnel_id")
	targetHostname := c.Param("hostname")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}
	if targetHostname == "" {
		utils.ErrorResponse(c, "Hostname is required", http.StatusBadRequest)
		return
	}

	// Create a simple struct to receive the JSON data
	var requestData struct {
		Hostname string `json:"hostname"`
		Service  string `json:"service"`
		Path     string `json:"path,omitempty"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.ErrorResponse(c, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if requestData.Hostname == "" {
		utils.ErrorResponse(c, "Hostname is required", http.StatusBadRequest)
		return
	}
	if requestData.Service == "" {
		utils.ErrorResponse(c, "Service is required", http.StatusBadRequest)
		return
	}

	// Create the hostname parameter using the helper function
	hostnameParam := models.NewPublicHostnameIngressParam(requestData.Hostname, requestData.Service, requestData.Path)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use the new function that updates both tunnel config and DNS record
	result, err := h.cfService.UpdateCloudflareTunnelPublicHostnameWithDNS(ctx, accountID, tunnelID, targetHostname, requestData.Hostname, hostnameParam)
	if err != nil {
		utils.ErrorResponse(c, "Failed to update public hostname: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":         "Public hostname updated successfully",
		"account_id":      accountID,
		"tunnel_id":       tunnelID,
		"target_hostname": targetHostname,
		"new_hostname":    requestData.Hostname,
		"config":          result,
	})
}

// DeletePublicHostname handles the DELETE /api/cloudflare/accounts/:accountId/tunnels/:tunnel_id/hostnames/:hostname endpoint
func (h *CloudflareTunnelHandler) DeletePublicHostname(c *gin.Context) {
	accountID := c.Param("accountId")
	tunnelID := c.Param("tunnel_id")
	targetHostname := c.Param("hostname")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}
	if targetHostname == "" {
		utils.ErrorResponse(c, "Hostname is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use the new function that deletes both tunnel config and DNS record
	result, err := h.cfService.DeleteCloudflareTunnelPublicHostnameWithDNS(ctx, accountID, tunnelID, targetHostname)
	if err != nil {
		utils.ErrorResponse(c, "Failed to delete public hostname: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":          "Public hostname deleted successfully",
		"account_id":       accountID,
		"tunnel_id":        tunnelID,
		"deleted_hostname": targetHostname,
		"config":           result,
	})
}
