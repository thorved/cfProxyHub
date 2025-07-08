package services

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareService struct {
	client *cloudflare.API
}

// NewCloudflareService creates a new Cloudflare service instance
func NewCloudflareService(apiToken, apiKey, email string) (*CloudflareService, error) {
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

	return &CloudflareService{
		client: client,
	}, nil
}
