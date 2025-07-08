package services

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

type CloudflareService struct {
	client *cloudflare.Client
}

// NewCloudflareService creates a new Cloudflare service instance
func NewCloudflareService(apiToken, apiKey, email string) (*CloudflareService, error) {
	var client *cloudflare.Client

	// Priority: API Token first, then API Key + Email
	if apiToken != "" && apiToken != "your_cloudflare_api_token_here" {
		client = cloudflare.NewClient(option.WithAPIToken(apiToken))
	} else if apiKey != "" && email != "" && apiKey != "your_cloudflare_api_key_here" {
		client = cloudflare.NewClient(
			option.WithAPIKey(apiKey),
			option.WithAPIEmail(email),
		)
	} else {
		return nil, fmt.Errorf("either API token or API key with email must be provided and properly configured")
	}

	service := &CloudflareService{
		client: client,
	}

	return service, nil
}
