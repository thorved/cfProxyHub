package models

import (
	"github.com/cloudflare/cloudflare-go/v4/accounts"
)

// Use v4 types directly - no conversion needed
type Account = accounts.Account
type AccountListParams = accounts.AccountListParams

// Response wrappers for API endpoints
type AccountResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    Account `json:"data"`
}

type AccountsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    []Account `json:"data"`
	Total   int       `json:"total"`
}

// Helper functions for creating request parameters
func NewAccountListParams() AccountListParams {
	return accounts.AccountListParams{}
}
