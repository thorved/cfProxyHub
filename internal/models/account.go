package models

import "time"

// AccountInfo represents a simplified view of a Cloudflare account
type AccountInfo struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type,omitempty"`
	Settings  *Settings  `json:"settings,omitempty"`
	CreatedOn *time.Time `json:"created_on,omitempty"`
}

// Settings represents account settings
type Settings struct {
	EnforceTwoFactor bool `json:"enforce_twofactor,omitempty"`
}

// AccountsResponse represents the response structure for accounts endpoint
type AccountsResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []AccountInfo `json:"data"`
	Total   int           `json:"total"`
}

// AccountResponse represents the response structure for single account endpoint
type AccountResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    AccountInfo `json:"data"`
}
