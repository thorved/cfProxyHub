package services

import (
	"context"
	"fmt"

	"cfPorxyHub/internal/models"

	"github.com/cloudflare/cloudflare-go/v4/accounts"
)

// GetAccounts retrieves all accounts associated with the API credentials
func (cs *CloudflareService) GetCloudflareAccounts(ctx context.Context) ([]models.Account, error) {
	autopager := cs.client.Accounts.ListAutoPaging(ctx, accounts.AccountListParams{})

	var result []models.Account
	for autopager.Next() {
		result = append(result, autopager.Current())
	}

	if autopager.Err() != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", autopager.Err())
	}

	return result, nil
}

// GetAccountByID retrieves a specific account by ID
func (cs *CloudflareService) GetCloudflareAccountByID(ctx context.Context, accountID string) (models.Account, error) {
	// The v3 API doesn't support getting a specific account by ID in the same way
	// We need to list all accounts and find the one with the matching ID
	autopager := cs.client.Accounts.ListAutoPaging(ctx, accounts.AccountListParams{})

	for autopager.Next() {
		account := autopager.Current()
		if account.ID == accountID {
			return account, nil
		}
	}

	if autopager.Err() != nil {
		return models.Account{}, fmt.Errorf("failed to fetch accounts: %w", autopager.Err())
	}

	return models.Account{}, fmt.Errorf("account with ID %s not found", accountID)
}
