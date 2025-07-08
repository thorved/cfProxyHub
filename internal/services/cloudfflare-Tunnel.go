package services

import (
	"context"
	"fmt"

	"cfPorxyHub/internal/models"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"
)

// GetCloudflareTunnels retrieves all tunnels for a specific account
func (cs *CloudflareService) GetCloudflareTunnels(ctx context.Context, accountID string) ([]models.TunnelListResponse, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	autopager := cs.client.ZeroTrust.Tunnels.Cloudflared.ListAutoPaging(ctx, zero_trust.TunnelCloudflaredListParams{
		AccountID: cloudflare.F(accountID),
	})

	var result []models.TunnelListResponse
	for autopager.Next() {
		tunnel := autopager.Current()
		result = append(result, tunnel)
	}

	if autopager.Err() != nil {
		return nil, fmt.Errorf("error listing tunnels for account %s: %w", accountID, autopager.Err())
	}

	return result, nil
}

// GetCloudflareTunnelByID retrieves a specific tunnel by ID for a specific account
func (cs *CloudflareService) GetCloudflareTunnelByID(ctx context.Context, accountID, tunnelID string) (models.Tunnel, error) {
	if accountID == "" {
		return models.Tunnel{}, fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return models.Tunnel{}, fmt.Errorf("tunnel ID is required")
	}

	tunnel, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Get(ctx, tunnelID, zero_trust.TunnelCloudflaredGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		return models.Tunnel{}, fmt.Errorf("error getting tunnel %s for account %s: %w", tunnelID, accountID, err)
	}

	return *tunnel, nil
}

// CreateCloudflareTunnel creates a new tunnel
func (cs *CloudflareService) CreateCloudflareTunnel(ctx context.Context, accountID string, request models.TunnelCreateRequest) (models.TunnelNewResponse, error) {
	if accountID == "" {
		return models.TunnelNewResponse{}, fmt.Errorf("account ID is required")
	}

	// Set the account ID in the request
	request.AccountID = cloudflare.F(accountID)

	createdTunnel, err := cs.client.ZeroTrust.Tunnels.Cloudflared.New(ctx, request)
	if err != nil {
		return models.TunnelNewResponse{}, fmt.Errorf("error creating tunnel for account %s: %w", accountID, err)
	}

	return *createdTunnel, nil
}

// UpdateCloudflareTunnel updates a specific tunnel
func (cs *CloudflareService) UpdateCloudflareTunnel(ctx context.Context, accountID, tunnelID string, request models.TunnelUpdateRequest) (models.TunnelEditResponse, error) {
	if accountID == "" {
		return models.TunnelEditResponse{}, fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return models.TunnelEditResponse{}, fmt.Errorf("tunnel ID is required")
	}

	// Set the account ID in the request
	request.AccountID = cloudflare.F(accountID)

	updatedTunnel, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Edit(ctx, tunnelID, request)
	if err != nil {
		return models.TunnelEditResponse{}, fmt.Errorf("error updating tunnel %s for account %s: %w", tunnelID, accountID, err)
	}

	return *updatedTunnel, nil
}

// DeleteCloudflareTunnel deletes a specific tunnel
func (cs *CloudflareService) DeleteCloudflareTunnel(ctx context.Context, accountID, tunnelID string) error {
	if accountID == "" {
		return fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return fmt.Errorf("tunnel ID is required")
	}

	_, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Delete(ctx, tunnelID, zero_trust.TunnelCloudflaredDeleteParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		return fmt.Errorf("error deleting tunnel %s for account %s: %w", tunnelID, accountID, err)
	}

	return nil
}

// GetCloudflareTunnelToken retrieves a token for a specific tunnel
func (cs *CloudflareService) GetCloudflareTunnelToken(ctx context.Context, accountID, tunnelID string) (string, error) {
	if accountID == "" {
		return "", fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return "", fmt.Errorf("tunnel ID is required")
	}

	tokenResponse, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Token.Get(ctx, tunnelID, zero_trust.TunnelCloudflaredTokenGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		return "", fmt.Errorf("error getting token for tunnel %s in account %s: %w", tunnelID, accountID, err)
	}

	// The token response might be an array or string, handle both cases
	if tokenResponse != nil {
		// For now, return a simple string representation
		return fmt.Sprintf("%v", tokenResponse), nil
	}

	return "", nil
}

// ListCloudflareTunnelsWithParams retrieves tunnels with specific parameters
func (cs *CloudflareService) ListCloudflareTunnelsWithParams(ctx context.Context, accountID string, params models.TunnelListParams) ([]models.TunnelListResponse, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	// Set the account ID in the params
	params.AccountID = cloudflare.F(accountID)

	autopager := cs.client.ZeroTrust.Tunnels.Cloudflared.ListAutoPaging(ctx, params)

	var result []models.TunnelListResponse
	for autopager.Next() {
		tunnel := autopager.Current()
		result = append(result, tunnel)
	}

	if autopager.Err() != nil {
		return nil, fmt.Errorf("error listing tunnels for account %s: %w", accountID, autopager.Err())
	}

	return result, nil
}
