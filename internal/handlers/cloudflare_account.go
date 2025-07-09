package handlers

import (
	"context"
	"net/http"
	"time"

	"cfProxyHub/internal/services"
	"cfProxyHub/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CloudflareAccountHandler struct {
	cfService *services.CloudflareService
}

// NewCloudflareAccountHandler creates a new Cloudflare account handler
func NewCloudflareAccountHandler(cfService *services.CloudflareService) *CloudflareAccountHandler {
	return &CloudflareAccountHandler{
		cfService: cfService,
	}
}

// GetAccounts handles the GET /api/cloudflare/accounts endpoint
func (h *CloudflareAccountHandler) GetAccounts(c *gin.Context) {
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
		"total":    len(accounts),
	})
}

// GetAccountByID handles the GET /api/cloudflare/accounts/:accountId endpoint
func (h *CloudflareAccountHandler) GetAccountByID(c *gin.Context) {
	accountID := c.Param("accountId")
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
