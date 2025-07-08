package handlers

import (
	"context"
	"net/http"
	"time"

	"cfPorxyHub/internal/services"
	"cfPorxyHub/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CloudflareHandler struct {
	cfService     *services.CloudflareService
	tunnelService *services.CloudflareTunnelService
}

// NewCloudflareHandler creates a new Cloudflare handler
func NewCloudflareHandler(cfService *services.CloudflareService, tunnelService *services.CloudflareTunnelService) *CloudflareHandler {
	return &CloudflareHandler{
		cfService:     cfService,
		tunnelService: tunnelService,
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

	tunnels, err := h.tunnelService.GetCloudflareTunnels(ctx, accountID)
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

	tunnel, err := h.tunnelService.GetCloudflareTunnelByID(ctx, accountID, tunnelID)
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
