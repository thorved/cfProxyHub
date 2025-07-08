package services

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// CurrentAccount represents the currently selected account
type CurrentAccount struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	SelectedAt   time.Time `json:"selected_at"`
	SessionToken string    `json:"-"` // Don't serialize session token
}

// AccountService manages current account selections
type AccountService struct {
	// In-memory storage for current account selections by session
	// In production, you would use Redis, database, or proper session storage
	sessions map[string]*CurrentAccount
	mutex    sync.RWMutex
}

// NewAccountService creates a new account service
func NewAccountService() *AccountService {
	return &AccountService{
		sessions: make(map[string]*CurrentAccount),
		mutex:    sync.RWMutex{},
	}
}

// GetCurrentAccount retrieves the current account for a session
func (s *AccountService) GetCurrentAccount(sessionToken string) (*CurrentAccount, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	account, exists := s.sessions[sessionToken]
	if !exists {
		return nil, errors.New("no account selected for this session")
	}

	// Check if selection is still valid (optional: add expiration logic)
	if time.Since(account.SelectedAt) > 24*time.Hour {
		// Selection is too old, remove it
		delete(s.sessions, sessionToken)
		return nil, errors.New("account selection expired")
	}

	return account, nil
}

// SetCurrentAccount sets the current account for a session
func (s *AccountService) SetCurrentAccount(sessionToken string, account *CurrentAccount) error {
	if sessionToken == "" {
		return errors.New("session token is required")
	}

	if account.ID == "" || account.Name == "" {
		return errors.New("account ID and name are required")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Set selection timestamp
	account.SelectedAt = time.Now()
	account.SessionToken = sessionToken

	// Store in session map
	s.sessions[sessionToken] = account

	return nil
}

// ClearCurrentAccount removes the current account selection for a session
func (s *AccountService) ClearCurrentAccount(sessionToken string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.sessions, sessionToken)
	return nil
}

// GetAllActiveSessions returns all active sessions with their account selections
// Useful for admin/debugging purposes
func (s *AccountService) GetAllActiveSessions() map[string]*CurrentAccount {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Create a copy to avoid external modifications
	result := make(map[string]*CurrentAccount)
	for sessionToken, account := range s.sessions {
		// Create a copy of the account
		accountCopy := &CurrentAccount{
			ID:         account.ID,
			Name:       account.Name,
			SelectedAt: account.SelectedAt,
		}
		result[sessionToken] = accountCopy
	}

	return result
}

// CleanupExpiredSessions removes expired account selections
// Should be called periodically (e.g., via a cron job)
func (s *AccountService) CleanupExpiredSessions() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	cleaned := 0
	now := time.Now()

	for sessionToken, account := range s.sessions {
		if now.Sub(account.SelectedAt) > 24*time.Hour {
			delete(s.sessions, sessionToken)
			cleaned++
		}
	}

	return cleaned
}

// ExportAccountSelection exports account selection to JSON (for persistence)
func (s *AccountService) ExportAccountSelection(sessionToken string) ([]byte, error) {
	account, err := s.GetCurrentAccount(sessionToken)
	if err != nil {
		return nil, err
	}

	return json.Marshal(account)
}

// ImportAccountSelection imports account selection from JSON (for persistence)
func (s *AccountService) ImportAccountSelection(sessionToken string, data []byte) error {
	var account CurrentAccount
	if err := json.Unmarshal(data, &account); err != nil {
		return err
	}

	return s.SetCurrentAccount(sessionToken, &account)
}
