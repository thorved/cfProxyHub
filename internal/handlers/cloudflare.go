package handlers

import (
	"context"
	"net/http"
	"time"

	"cfPorxyHub/internal/models"
	"cfPorxyHub/internal/services"
	"cfPorxyHub/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CloudflareHandler struct {
	cfService *services.CloudflareService
}

// NewCloudflareHandler creates a new Cloudflare handler
func NewCloudflareHandler(cfService *services.CloudflareService) *CloudflareHandler {
	return &CloudflareHandler{
		cfService: cfService,
	}
}

// GetAccounts handles the GET /api/cloudflare/accounts endpoint
func (h *CloudflareHandler) GetAccounts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accounts, err := h.cfService.GetCloudflareAccounts(ctx)
	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":  "Accounts retrieved successfully",
		"accounts": accounts,
	})
}

// GetAccountByID handles the GET /api/cloudflare/accounts/:id endpoint
func (h *CloudflareHandler) GetAccountByID(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	account, err := h.cfService.GetCloudflareAccountByID(ctx, accountID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Account retrieved successfully",
		"account": account,
	})
}

// GetTunnelsByAccountID handles the GET /api/cloudflare/accounts/:id/tunnels endpoint
func (h *CloudflareHandler) GetTunnelsByAccountID(c *gin.Context) {
	accountID := c.Param("id")
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

// GetTunnelByID handles the GET /api/cloudflare/accounts/:id/tunnels/:tunnel_id endpoint
func (h *CloudflareHandler) GetTunnelByID(c *gin.Context) {
	accountID := c.Param("id")
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

// CreateTunnel handles the POST /api/cloudflare/accounts/:id/tunnels endpoint
func (h *CloudflareHandler) CreateTunnel(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}

	var request models.TunnelCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

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

// UpdateTunnel handles the PUT /api/cloudflare/accounts/:id/tunnels/:tunnel_id endpoint
func (h *CloudflareHandler) UpdateTunnel(c *gin.Context) {
	accountID := c.Param("id")
	tunnelID := c.Param("tunnel_id")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if tunnelID == "" {
		utils.ErrorResponse(c, "Tunnel ID is required", http.StatusBadRequest)
		return
	}

	var request models.TunnelUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

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

// DeleteTunnel handles the DELETE /api/cloudflare/accounts/:id/tunnels/:tunnel_id endpoint
func (h *CloudflareHandler) DeleteTunnel(c *gin.Context) {
	accountID := c.Param("id")
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

// GetTunnelToken handles the GET /api/cloudflare/accounts/:id/tunnels/:tunnel_id/token endpoint
func (h *CloudflareHandler) GetTunnelToken(c *gin.Context) {
	accountID := c.Param("id")
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
