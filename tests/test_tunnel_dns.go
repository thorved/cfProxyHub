package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cfProxyHub/internal/models"
	"cfProxyHub/internal/services"
)

// This script tests the automatic DNS record creation when adding a public hostname to a tunnel
func main() {
	// Replace these with real values for testing
	apiToken := "your-api-token"
	accountID := "your-account-id"
	tunnelID := "your-tunnel-id"
	hostname := "test.example.com" // Replace with a domain you own in Cloudflare

	// Create our service
	cfService, err := services.NewCloudflareService(apiToken, "", "")
	if err != nil {
		log.Fatalf("Error creating Cloudflare service: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create the hostname parameter
	hostnameParam := models.NewPublicHostnameIngressParam(hostname, "https://localhost:8080", "")

	// Call the function we want to test
	fmt.Printf("Adding hostname %s to tunnel %s in account %s...\n", hostname, tunnelID, accountID)
	result, err := cfService.CreateCloudflareTunnelPublicHostnameWithDNS(ctx, accountID, tunnelID, hostname, hostnameParam)
	if err != nil {
		log.Fatalf("Error creating public hostname with DNS: %v", err)
	}
	fmt.Printf("Successfully added hostname: %+v\n", result)

	// Verify DNS record was created
	// Extract domain from hostname
	domain := extractDomain(hostname)

	// Find the zone for this domain
	fmt.Printf("Looking up zone for domain %s...\n", domain)
	zone, err := cfService.GetZoneByName(ctx, accountID, domain)
	if err != nil {
		log.Fatalf("Error getting zone: %v", err)
	}
	fmt.Printf("Found zone: %s (ID: %s)\n", zone.Name, zone.ID)

	// Check for DNS record
	fmt.Printf("Checking for DNS record for hostname %s...\n", hostname)
	records, err := cfService.GetDNSRecordsByName(ctx, zone.ID, hostname)
	if err != nil {
		log.Fatalf("Error getting DNS records: %v", err)
	}

	if len(records) == 0 {
		log.Fatalf("No DNS record found for hostname %s", hostname)
	}

	// Display the record
	for _, record := range records {
		fmt.Printf("DNS Record: %s (%s) -> %s\n", record.Name, record.Type, record.Content)
		// Verify it points to the tunnel
		expectedTarget := fmt.Sprintf("%s.cfargotunnel.com", tunnelID)
		if record.Type == "CNAME" && record.Content == expectedTarget {
			fmt.Printf("✅ SUCCESS: DNS record correctly points to %s\n", expectedTarget)
		}
	}
}

// Helper function to extract domain from hostname
func extractDomain(hostname string) string {
	parts := strings.Split(hostname, ".")
	if len(parts) < 2 {
		return hostname
	}
	// Return the last two parts joined (this is a simple approach)
	return strings.Join(parts[len(parts)-2:], ".")
}
