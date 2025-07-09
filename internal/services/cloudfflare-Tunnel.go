package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cfPorxyHub/internal/models"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"
)

// GetCloudflareTunnels retrieves all active (non-deleted) tunnels for a specific account
func (cs *CloudflareService) GetCloudflareTunnels(ctx context.Context, accountID string) ([]models.TunnelListResponse, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	autopager := cs.client.ZeroTrust.Tunnels.Cloudflared.ListAutoPaging(ctx, zero_trust.TunnelCloudflaredListParams{
		AccountID: cloudflare.F(accountID),
		IsDeleted: cloudflare.F(false),
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

	// Always set is_deleted to false to only get active tunnels
	params.IsDeleted = cloudflare.F(false)

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

// Public Hostname Management Functions

// GetCloudflareTunnelPublicHostnames retrieves all public hostnames (ingress rules) for a specific tunnel
func (cs *CloudflareService) GetCloudflareTunnelPublicHostnames(ctx context.Context, accountID, tunnelID string) ([]models.PublicHostnameIngress, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return nil, fmt.Errorf("tunnel ID is required")
	}

	// Get tunnel configuration which contains the ingress rules (public hostnames)
	config, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Get(ctx, tunnelID, zero_trust.TunnelCloudflaredConfigurationGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting tunnel configuration for tunnel %s in account %s: %w", tunnelID, accountID, err)
	}

	// Extract ingress rules (public hostnames) from the configuration
	if config.Config.Ingress == nil {
		return []models.PublicHostnameIngress{}, nil
	}

	return config.Config.Ingress, nil
}

// CreateCloudflareTunnelPublicHostname creates a new public hostname (ingress rule) for a specific tunnel
func (cs *CloudflareService) CreateCloudflareTunnelPublicHostname(ctx context.Context, accountID, tunnelID string, hostname models.PublicHostnameIngressParam) (models.TunnelConfigurationUpdateResponse, error) {
	if accountID == "" {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("tunnel ID is required")
	}

	// TODO: Add domain validation once we figure out the field access pattern
	// For now, we'll rely on Cloudflare's built-in validation
	//
	// Example of what we want to do:
	// hostnameValue := hostname.Hostname.ValueOrZero()
	// if hostnameValue != "" {
	//     if err := cs.validateDomainOwnership(ctx, accountID, hostnameValue); err != nil {
	//         return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("domain validation failed: %w", err)
	//     }
	// }

	// TODO: Add domain validation for updated hostname once we understand field access pattern
	// For now, we'll rely on Cloudflare's built-in validation

	// First, get the current configuration
	currentConfig, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Get(ctx, tunnelID, zero_trust.TunnelCloudflaredConfigurationGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		// If configuration doesn't exist, start with an empty one
		fmt.Printf("Warning: Could not get existing configuration (this is normal for new tunnels): %v\n", err)
		currentConfig = &zero_trust.TunnelCloudflaredConfigurationGetResponse{
			Config: zero_trust.TunnelCloudflaredConfigurationGetResponseConfig{},
		}
	}

	// Check if hostname already exists
	if currentConfig.Config.Ingress != nil {
		// TODO: Fix hostname comparison once we understand the field access pattern
		// For now, we'll skip the duplicate check and rely on Cloudflare's validation
	}

	// Create the new ingress rules array with the new hostname
	var newIngress []models.PublicHostnameIngressParam

	// Add existing ingress rules (excluding any catch-all rules)
	if currentConfig != nil && currentConfig.Config.Ingress != nil {
		for _, existing := range currentConfig.Config.Ingress {
			// Skip catch-all rules (rules without hostname)
			if existing.Hostname != "" {
				existingRule := models.PublicHostnameIngressParam{
					Hostname:      cloudflare.F(existing.Hostname),
					Service:       cloudflare.F(existing.Service),
					OriginRequest: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngressOriginRequest{}),
				}
				if existing.Path != "" {
					existingRule.Path = cloudflare.F(existing.Path)
				}
				newIngress = append(newIngress, existingRule)
			}
		}
	}

	// Add the new hostname
	newIngress = append(newIngress, hostname)

	// Always add a catch-all rule at the end (required by Cloudflare)
	// This rule will return a 404 for any hostname that doesn't match the rules above
	catchAllRule := models.PublicHostnameIngressParam{
		Service: cloudflare.F("http_status:404"),
	}
	newIngress = append(newIngress, catchAllRule)

	// Update the tunnel configuration
	updateParams := zero_trust.TunnelCloudflaredConfigurationUpdateParams{
		AccountID: cloudflare.F(accountID),
		Config: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfig{
			Ingress: cloudflare.F(newIngress),
		}),
	}

	updatedConfig, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Update(ctx, tunnelID, updateParams)
	if err != nil {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("error creating public hostname for tunnel %s in account %s: %w", tunnelID, accountID, err)
	}

	return *updatedConfig, nil
}

// UpdateCloudflareTunnelPublicHostname updates a specific public hostname (ingress rule) for a tunnel
func (cs *CloudflareService) UpdateCloudflareTunnelPublicHostname(ctx context.Context, accountID, tunnelID, targetHostname string, updatedHostname models.PublicHostnameIngressParam) (models.TunnelConfigurationUpdateResponse, error) {
	log.Printf("UpdateCloudflareTunnelPublicHostname called with accountID: %s, tunnelID: %s, targetHostname: %s", accountID, tunnelID, targetHostname)

	if accountID == "" {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("tunnel ID is required")
	}
	if targetHostname == "" {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("target hostname is required")
	}

	log.Printf("Getting current tunnel configuration...")

	// Get the current configuration
	currentConfig, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Get(ctx, tunnelID, zero_trust.TunnelCloudflaredConfigurationGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		log.Printf("Error getting current tunnel configuration: %v", err)
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("error getting current tunnel configuration for tunnel %s in account %s: %w", tunnelID, accountID, err)
	}

	log.Printf("Current configuration retrieved successfully, processing ingress rules...")

	// Create the new ingress rules array with the updated hostname
	var newIngress []models.PublicHostnameIngressParam
	found := false

	if currentConfig.Config.Ingress != nil {
		log.Printf("Found %d existing ingress rules", len(currentConfig.Config.Ingress))
		for i, existing := range currentConfig.Config.Ingress {
			log.Printf("Rule %d: hostname=%s, service=%s", i, existing.Hostname, existing.Service)

			// Skip catch-all rules (rules without hostname)
			if existing.Hostname == "" {
				log.Printf("Skipping catch-all rule %d", i)
				continue
			}

			if existing.Hostname == targetHostname {
				log.Printf("Found target hostname %s at rule %d, replacing with updated hostname", targetHostname, i)
				// Replace with the updated hostname
				newIngress = append(newIngress, updatedHostname)
				found = true
			} else {
				log.Printf("Keeping existing rule %d: %s", i, existing.Hostname)
				// Keep existing
				newIngress = append(newIngress, models.PublicHostnameIngressParam{
					Hostname:      cloudflare.F(existing.Hostname),
					Service:       cloudflare.F(existing.Service),
					Path:          cloudflare.F(existing.Path),
					OriginRequest: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngressOriginRequest{}),
				})
			}
		}
	} else {
		log.Printf("No existing ingress rules found")
	}

	if !found {
		log.Printf("Target hostname %s not found in tunnel %s", targetHostname, tunnelID)
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("hostname %s not found in tunnel %s", targetHostname, tunnelID)
	}

	log.Printf("Successfully found and updated target hostname, adding catch-all rule")

	// Always add a catch-all rule at the end (required by Cloudflare)
	catchAllRule := models.PublicHostnameIngressParam{
		Service: cloudflare.F("http_status:404"),
	}
	newIngress = append(newIngress, catchAllRule)

	log.Printf("Updating tunnel configuration with %d total rules", len(newIngress))

	// Update the tunnel configuration
	updateParams := zero_trust.TunnelCloudflaredConfigurationUpdateParams{
		AccountID: cloudflare.F(accountID),
		Config: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfig{
			Ingress: cloudflare.F(newIngress),
		}),
	}

	updatedConfig, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Update(ctx, tunnelID, updateParams)
	if err != nil {
		log.Printf("Error updating tunnel configuration: %v", err)
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("error updating public hostname for tunnel %s in account %s: %w", tunnelID, accountID, err)
	}

	log.Printf("Successfully updated tunnel configuration")
	return *updatedConfig, nil
}

// DeleteCloudflareTunnelPublicHostname deletes a specific public hostname (ingress rule) from a tunnel
func (cs *CloudflareService) DeleteCloudflareTunnelPublicHostname(ctx context.Context, accountID, tunnelID, targetHostname string) (models.TunnelConfigurationUpdateResponse, error) {
	if accountID == "" {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("account ID is required")
	}
	if tunnelID == "" {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("tunnel ID is required")
	}
	if targetHostname == "" {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("target hostname is required")
	}

	// Get the current configuration
	currentConfig, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Get(ctx, tunnelID, zero_trust.TunnelCloudflaredConfigurationGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("error getting current tunnel configuration for tunnel %s in account %s: %w", tunnelID, accountID, err)
	}

	// Create the new ingress rules array without the target hostname
	var newIngress []models.PublicHostnameIngressParam
	found := false

	if currentConfig.Config.Ingress != nil {
		for _, existing := range currentConfig.Config.Ingress {
			// Skip catch-all rules (rules without hostname)
			if existing.Hostname == "" {
				continue
			}

			if existing.Hostname != targetHostname {
				// Keep this hostname
				newIngress = append(newIngress, models.PublicHostnameIngressParam{
					Hostname:      cloudflare.F(existing.Hostname),
					Service:       cloudflare.F(existing.Service),
					Path:          cloudflare.F(existing.Path),
					OriginRequest: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngressOriginRequest{}),
				})
			} else {
				found = true
			}
		}
	}

	if !found {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("hostname %s not found in tunnel %s", targetHostname, tunnelID)
	}

	// Always add a catch-all rule at the end (required by Cloudflare)
	catchAllRule := models.PublicHostnameIngressParam{
		Service: cloudflare.F("http_status:404"),
	}
	newIngress = append(newIngress, catchAllRule)

	// Update the tunnel configuration
	updateParams := zero_trust.TunnelCloudflaredConfigurationUpdateParams{
		AccountID: cloudflare.F(accountID),
		Config: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfig{
			Ingress: cloudflare.F(newIngress),
		}),
	}

	updatedConfig, err := cs.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Update(ctx, tunnelID, updateParams)
	if err != nil {
		return models.TunnelConfigurationUpdateResponse{}, fmt.Errorf("error deleting public hostname for tunnel %s in account %s: %w", tunnelID, accountID, err)
	}

	return *updatedConfig, nil
}

// CreateCloudflareTunnelPublicHostnameWithDNS creates a new public hostname (ingress rule) for a specific tunnel and automatically creates DNS record
func (cs *CloudflareService) CreateCloudflareTunnelPublicHostnameWithDNS(ctx context.Context, accountID, tunnelID, hostnameValue string, hostname models.PublicHostnameIngressParam) (models.TunnelConfigurationUpdateResponse, error) {
	// First create the tunnel configuration
	updatedConfig, err := cs.CreateCloudflareTunnelPublicHostname(ctx, accountID, tunnelID, hostname)
	if err != nil {
		return models.TunnelConfigurationUpdateResponse{}, err
	}

	// Try to create DNS record for the hostname
	if hostnameValue != "" {
		// Extract domain from hostname
		domain := extractDomain(hostnameValue)

		// Find the zone for this domain
		zone, err := cs.GetZoneByName(ctx, accountID, domain)
		if err != nil {
			// Log warning but don't fail the tunnel configuration
			log.Printf("Warning: Could not find zone for domain %s: %v", domain, err)
		} else {
			// Create CNAME record pointing to the tunnel
			// For tunnel hostnames, we typically want them proxied through Cloudflare for security benefits
			_, err = cs.CreateTunnelCNAMERecord(ctx, zone.ID, hostnameValue, tunnelID, true)
			if err != nil {
				// Log warning but don't fail the tunnel configuration
				log.Printf("Warning: Could not create DNS record for hostname %s: %v", hostnameValue, err)
			} else {
				log.Printf("Successfully created DNS record for hostname %s -> %s.cfargotunnel.com", hostnameValue, tunnelID)
			}
		}
	}

	return updatedConfig, nil
}

// UpdateCloudflareTunnelPublicHostnameWithDNS updates a public hostname (ingress rule) and its DNS record
func (cs *CloudflareService) UpdateCloudflareTunnelPublicHostnameWithDNS(ctx context.Context, accountID, tunnelID, targetHostname, newHostnameValue string, updatedHostname models.PublicHostnameIngressParam) (models.TunnelConfigurationUpdateResponse, error) {
	log.Printf("UpdateCloudflareTunnelPublicHostnameWithDNS called with accountID: %s, tunnelID: %s, targetHostname: %s, newHostnameValue: %s", accountID, tunnelID, targetHostname, newHostnameValue)

	// First update the tunnel configuration
	updatedConfig, err := cs.UpdateCloudflareTunnelPublicHostname(ctx, accountID, tunnelID, targetHostname, updatedHostname)
	if err != nil {
		log.Printf("Error updating tunnel configuration: %v", err)
		return models.TunnelConfigurationUpdateResponse{}, err
	}

	log.Printf("Successfully updated tunnel configuration")

	// If the hostname value changed, we need to update DNS records
	if targetHostname != newHostnameValue && newHostnameValue != "" {
		// Extract domains
		oldDomain := extractDomain(targetHostname)
		newDomain := extractDomain(newHostnameValue)

		// Find the zone for the old domain (to delete the old record)
		if oldZone, err := cs.GetZoneByName(ctx, accountID, oldDomain); err == nil {
			// Try to find and delete the old DNS record
			oldRecords, err := cs.GetDNSRecordsByName(ctx, oldZone.ID, targetHostname)
			if err == nil && len(oldRecords) > 0 {
				for _, record := range oldRecords {
					if record.Type == "CNAME" && strings.Contains(record.Content, "cfargotunnel.com") {
						// Delete the old record
						err := cs.DeleteDNSRecord(ctx, oldZone.ID, record.ID)
						if err != nil {
							log.Printf("Warning: Could not delete old DNS record for %s: %v", targetHostname, err)
						} else {
							log.Printf("Successfully deleted old DNS record for %s", targetHostname)
						}
					}
				}
			}
		}

		// Create a new DNS record for the new hostname
		newZone, err := cs.GetZoneByName(ctx, accountID, newDomain)
		if err != nil {
			log.Printf("Warning: Could not find zone for new domain %s: %v", newDomain, err)
		} else {
			// Create CNAME record pointing to the tunnel
			// For tunnel hostnames, we typically want them proxied through Cloudflare for security benefits
			_, err = cs.CreateTunnelCNAMERecord(ctx, newZone.ID, newHostnameValue, tunnelID, true)
			if err != nil {
				log.Printf("Warning: Could not create DNS record for new hostname %s: %v", newHostnameValue, err)
			} else {
				log.Printf("Successfully created DNS record for hostname %s -> %s.cfargotunnel.com", newHostnameValue, tunnelID)
			}
		}
	}

	return updatedConfig, nil
}

// DeleteCloudflareTunnelPublicHostnameWithDNS deletes a public hostname (ingress rule) and its DNS record
func (cs *CloudflareService) DeleteCloudflareTunnelPublicHostnameWithDNS(ctx context.Context, accountID, tunnelID, targetHostname string) (models.TunnelConfigurationUpdateResponse, error) {
	// First delete the tunnel configuration
	updatedConfig, err := cs.DeleteCloudflareTunnelPublicHostname(ctx, accountID, tunnelID, targetHostname)
	if err != nil {
		return models.TunnelConfigurationUpdateResponse{}, err
	}

	// Try to delete the DNS record as well
	if targetHostname != "" {
		// Extract domain from hostname
		domain := extractDomain(targetHostname)

		// Find the zone for this domain
		zone, err := cs.GetZoneByName(ctx, accountID, domain)
		if err != nil {
			log.Printf("Warning: Could not find zone for domain %s: %v", domain, err)
		} else {
			// Try to find and delete the DNS record
			records, err := cs.GetDNSRecordsByName(ctx, zone.ID, targetHostname)
			if err != nil {
				log.Printf("Warning: Could not get DNS records for hostname %s: %v", targetHostname, err)
			} else if len(records) > 0 {
				for _, record := range records {
					if record.Type == "CNAME" && strings.Contains(record.Content, "cfargotunnel.com") {
						err := cs.DeleteDNSRecord(ctx, zone.ID, record.ID)
						if err != nil {
							log.Printf("Warning: Could not delete DNS record for hostname %s: %v", targetHostname, err)
						} else {
							log.Printf("Successfully deleted DNS record for hostname %s", targetHostname)
						}
					}
				}
			}
		}
	}

	return updatedConfig, nil
}

// CreateTunnelCNAMERecord creates a CNAME record for a hostname pointing to a tunnel
func (cs *CloudflareService) CreateTunnelCNAMERecord(ctx context.Context, zoneID, hostname, tunnelID string, proxied bool) (*models.DNSRecord, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("zone ID is required")
	}
	if hostname == "" {
		return nil, fmt.Errorf("hostname is required")
	}
	if tunnelID == "" {
		return nil, fmt.Errorf("tunnel ID is required")
	}

	// Construct the tunnel domain with the Cloudflare-specific suffix
	tunnelDomain := fmt.Sprintf("%s.cfargotunnel.com", tunnelID)

	log.Printf("Creating CNAME record for tunnel: %s -> %s (proxied: %t)", hostname, tunnelDomain, proxied)

	// Use the generic CNAME record creation function
	return cs.CreateCNAMERecord(ctx, zoneID, hostname, tunnelDomain, proxied)
}

// validateDomainOwnership checks if the given hostname's domain is owned by the account
func (cs *CloudflareService) validateDomainOwnership(ctx context.Context, accountID, hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	// Extract domain from hostname
	parts := strings.Split(hostname, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid hostname format")
	}

	// Find the domain part (assume it's the last two parts for now)
	// This is a simple approach - in reality, you might need more sophisticated logic
	domain := strings.Join(parts[len(parts)-2:], ".")

	// Get zones for the account to check if domain is owned
	zones, err := cs.GetActiveZones(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to verify domain ownership: %w", err)
	}

	// Check if any zone matches the domain
	for _, zone := range zones {
		if zone.Name == domain {
			return nil // Domain is owned by this account
		}
	}

	return fmt.Errorf("domain %s is not owned by account %s", domain, accountID)
}
