package services

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareTunnelService struct {
	client *cloudflare.API
}

// NewCloudflareTunnelService creates a new Cloudflare tunnel service instance
func NewCloudflareTunnelService(apiToken, apiKey, email string) (*CloudflareTunnelService, error) {
	var client *cloudflare.API
	var err error

	// Prefer API Token over API Key + Email
	if apiToken != "" {
		client, err = cloudflare.NewWithAPIToken(apiToken)
	} else if apiKey != "" && email != "" {
		client, err = cloudflare.New(apiKey, email)
	} else {
		return nil, fmt.Errorf("either API token or API key with email must be provided")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudflare client: %w", err)
	}

	return &CloudflareTunnelService{
		client: client,
	}, nil
}

func (cs *CloudflareTunnelService) GetCloudflareTunnels(ctx context.Context, accountID string) ([]cloudflare.Tunnel, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	tunnels, _, err := cs.client.ListTunnels(ctx, cloudflare.AccountIdentifier(accountID), cloudflare.TunnelListParams{})
	if err != nil {
		return nil, fmt.Errorf("error listing tunnels for account %s: %w", accountID, err)
	}
	return tunnels, nil
}

func (cs *CloudflareTunnelService) GetCloudflareTunnelByID(ctx context.Context, accountID, tunnelID string) (cloudflare.Tunnel, error) {
	if accountID == "" {
		return cloudflare.Tunnel{}, fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return cloudflare.Tunnel{}, fmt.Errorf("tunnel ID is required")
	}

	tunnel, err := cs.client.GetTunnel(ctx, cloudflare.AccountIdentifier(accountID), tunnelID)
	if err != nil {
		return cloudflare.Tunnel{}, fmt.Errorf("error getting tunnel %s for account %s: %w", tunnelID, accountID, err)
	}
	return tunnel, nil
}
