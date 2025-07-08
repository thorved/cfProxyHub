package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

// GetAccountsHTML handles HTMX requests for accounts and returns HTML fragments
func (h *CloudflareHandler) GetAccountsHTML(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accounts, err := h.cfService.GetCloudflareAccounts(ctx)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "", gin.H{
			"error": "Failed to fetch accounts: " + err.Error(),
		})
		return
	}

	if len(accounts) == 0 {
		c.String(http.StatusOK, `<div class="loading">No accounts found.</div>`)
		return
	}

	var html string
	for _, account := range accounts {
		createdDate := ""
		if !account.CreatedOn.IsZero() {
			createdDate = account.CreatedOn.Format("2006-01-02")
		}

		accountType := ""
		if account.Type != "" {
			accountType = fmt.Sprintf(`<div class="account-type"><strong>Type:</strong> %s</div>`, account.Type)
		}

		createdOnDiv := ""
		if createdDate != "" {
			createdOnDiv = fmt.Sprintf(`<div class="account-type"><strong>Created:</strong> %s</div>`, createdDate)
		}

		html += fmt.Sprintf(`
			<div class="account-card">
				<div class="account-name">%s</div>
				<div><strong>ID:</strong> <span class="account-id">%s</span></div>
				%s
				%s
			</div>
		`, account.Name, account.ID, accountType, createdOnDiv)
	}

	c.String(http.StatusOK, html)
}
