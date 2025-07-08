package services

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

// GetAccounts retrieves all accounts associated with the API credentials
func (cs *CloudflareService) GetCloudflareAccounts(ctx context.Context) ([]cloudflare.Account, error) {
	accounts, _, err := cs.client.Accounts(ctx, cloudflare.AccountsListParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	return accounts, nil
}

// GetAccountByID retrieves a specific account by ID
func (cs *CloudflareService) GetCloudflareAccountByID(ctx context.Context, accountID string) (cloudflare.Account, error) {
	account, _, err := cs.client.Account(ctx, accountID)
	if err != nil {
		return cloudflare.Account{}, fmt.Errorf("failed to fetch account %s: %w", accountID, err)
	}

	return account, nil
}
