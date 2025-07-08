package handlers

import (
	"cfPorxyHub/internal/services"
	"cfPorxyHub/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountHandler handles account-related operations
type AccountHandler struct {
	accountService *services.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// GetCurrentAccount handles GET /api/current-account
func (h *AccountHandler) GetCurrentAccount(c *gin.Context) {
	// Get session token to identify user
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		utils.ErrorResponse(c, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get current account from service
	account, err := h.accountService.GetCurrentAccount(sessionToken)
	if err != nil {
		// No account selected server-side, client should check localStorage
		utils.SuccessResponse(c, gin.H{
			"success": false,
			"message": "No account selected server-side, check localStorage",
			"account": nil,
		})
		return
	}

	utils.SuccessResponse(c, gin.H{
		"success": true,
		"message": "Current account retrieved",
		"account": account,
	})
}

// SetCurrentAccount handles POST /api/current-account
func (h *AccountHandler) SetCurrentAccount(c *gin.Context) {
	var accountData struct {
		ID   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&accountData); err != nil {
		utils.ErrorResponse(c, "Invalid account data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get session token to identify user
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		utils.ErrorResponse(c, "Session not found", http.StatusUnauthorized)
		return
	}

	// Set current account via service
	account := &services.CurrentAccount{
		ID:   accountData.ID,
		Name: accountData.Name,
	}

	err = h.accountService.SetCurrentAccount(sessionToken, account)
	if err != nil {
		utils.ErrorResponse(c, "Failed to save account selection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"success": true,
		"message": "Account selection saved",
		"account": gin.H{
			"id":   account.ID,
			"name": account.Name,
		},
	})
}

// ClearCurrentAccount handles DELETE /api/current-account
func (h *AccountHandler) ClearCurrentAccount(c *gin.Context) {
	// Get session token to identify user
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		utils.ErrorResponse(c, "Session not found", http.StatusUnauthorized)
		return
	}

	// Clear current account via service
	err = h.accountService.ClearCurrentAccount(sessionToken)
	if err != nil {
		utils.ErrorResponse(c, "Failed to clear account selection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Account selection cleared",
	})
}
