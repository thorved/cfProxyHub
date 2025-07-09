package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cfPorxyHub/internal/models"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

// GetZonesByAccountID retrieves all zones for a specific account
func (s *CloudflareService) GetZonesByAccountID(ctx context.Context, accountID string) ([]models.Zone, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching zones for account: %s", accountID)

	params := zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{
			ID: cloudflare.F(accountID),
		}),
	}

	zoneList, err := s.client.Zones.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch zones: %w", err)
	}

	var result []models.Zone
	for _, zone := range zoneList.Result {
		result = append(result, models.Zone{
			ID:                  zone.ID,
			Name:                zone.Name,
			Status:              string(zone.Status),
			Type:                string(zone.Type),
			DevelopmentMode:     int(zone.DevelopmentMode),
			NameServers:         zone.NameServers,
			OriginalNameServers: zone.OriginalNameServers,
			OriginalRegistrar:   zone.OriginalRegistrar,
			OriginalDNSHost:     zone.OriginalDnshost,
			ModifiedOn:          zone.ModifiedOn,
			CreatedOn:           zone.CreatedOn,
			ActivatedOn:         zone.ActivatedOn,
		})
	}

	log.Printf("Successfully fetched %d zones for account %s", len(result), accountID)
	return result, nil
}

// GetZoneByID retrieves a specific zone by ID
func (s *CloudflareService) GetZoneByID(ctx context.Context, zoneID string) (*models.Zone, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("zone ID is required")
	}

	log.Printf("Fetching zone: %s", zoneID)

	zone, err := s.client.Zones.Get(ctx, zones.ZoneGetParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch zone: %w", err)
	}

	result := models.Zone{
		ID:                  zone.ID,
		Name:                zone.Name,
		Status:              string(zone.Status),
		Type:                string(zone.Type),
		DevelopmentMode:     int(zone.DevelopmentMode),
		NameServers:         zone.NameServers,
		OriginalNameServers: zone.OriginalNameServers,
		OriginalRegistrar:   zone.OriginalRegistrar,
		OriginalDNSHost:     zone.OriginalDnshost,
		ModifiedOn:          zone.ModifiedOn,
		CreatedOn:           zone.CreatedOn,
		ActivatedOn:         zone.ActivatedOn,
	}
	log.Printf("Successfully fetched zone: %s (%s)", zone.Name, zone.ID)
	return &result, nil
}

// GetZoneByName retrieves a zone by domain name
func (s *CloudflareService) GetZoneByName(ctx context.Context, accountID, domainName string) (*models.Zone, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}
	if domainName == "" {
		return nil, fmt.Errorf("domain name is required")
	}

	log.Printf("Fetching zone by name: %s for account: %s", domainName, accountID)

	params := zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{
			ID: cloudflare.F(accountID),
		}),
		Name: cloudflare.F(domainName),
	}

	zoneList, err := s.client.Zones.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch zone by name: %w", err)
	}

	if len(zoneList.Result) == 0 {
		return nil, fmt.Errorf("zone not found: %s", domainName)
	}

	zone := zoneList.Result[0]
	result := models.Zone{
		ID:                  zone.ID,
		Name:                zone.Name,
		Status:              string(zone.Status),
		Type:                string(zone.Type),
		DevelopmentMode:     int(zone.DevelopmentMode),
		NameServers:         zone.NameServers,
		OriginalNameServers: zone.OriginalNameServers,
		OriginalRegistrar:   zone.OriginalRegistrar,
		OriginalDNSHost:     zone.OriginalDnshost,
		ModifiedOn:          zone.ModifiedOn,
		CreatedOn:           zone.CreatedOn,
		ActivatedOn:         zone.ActivatedOn,
	}
	log.Printf("Successfully fetched zone by name: %s (%s)", result.Name, result.ID)
	return &result, nil
}

// GetActiveZones retrieves only active zones for an account
func (s *CloudflareService) GetActiveZones(ctx context.Context, accountID string) ([]models.Zone, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching active zones for account: %s", accountID)

	params := zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{
			ID: cloudflare.F(accountID),
		}),
		Status: cloudflare.F(zones.ZoneListParamsStatusActive),
	}

	zoneList, err := s.client.Zones.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active zones: %w", err)
	}

	var result []models.Zone
	for _, zone := range zoneList.Result {
		result = append(result, models.Zone{
			ID:                  zone.ID,
			Name:                zone.Name,
			Status:              string(zone.Status),
			Type:                string(zone.Type),
			DevelopmentMode:     int(zone.DevelopmentMode),
			NameServers:         zone.NameServers,
			OriginalNameServers: zone.OriginalNameServers,
			OriginalRegistrar:   zone.OriginalRegistrar,
			OriginalDNSHost:     zone.OriginalDnshost,
			ModifiedOn:          zone.ModifiedOn,
			CreatedOn:           zone.CreatedOn,
			ActivatedOn:         zone.ActivatedOn,
		})
	}

	log.Printf("Successfully fetched %d active zones for account %s", len(result), accountID)
	return result, nil
}

// SearchZones searches for zones by name pattern
func (s *CloudflareService) SearchZones(ctx context.Context, accountID, searchTerm string) ([]models.Zone, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Searching zones for account: %s, term: %s", accountID, searchTerm)

	params := zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{
			ID: cloudflare.F(accountID),
		}),
	}

	// If search term is provided, add it to the search
	if searchTerm != "" {
		params.Name = cloudflare.F(searchTerm)
	}

	zoneList, err := s.client.Zones.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search zones: %w", err)
	}

	var result []models.Zone
	for _, zone := range zoneList.Result {
		result = append(result, models.Zone{
			ID:                  zone.ID,
			Name:                zone.Name,
			Status:              string(zone.Status),
			Type:                string(zone.Type),
			DevelopmentMode:     int(zone.DevelopmentMode),
			NameServers:         zone.NameServers,
			OriginalNameServers: zone.OriginalNameServers,
			OriginalRegistrar:   zone.OriginalRegistrar,
			OriginalDNSHost:     zone.OriginalDnshost,
			ModifiedOn:          zone.ModifiedOn,
			CreatedOn:           zone.CreatedOn,
			ActivatedOn:         zone.ActivatedOn,
		})
	}

	log.Printf("Successfully found %d zones matching search term '%s' for account %s", len(result), searchTerm, accountID)
	return result, nil
}

// DNS Record Management Functions

// CreateDNSRecord creates a new DNS record for a zone
func (s *CloudflareService) CreateDNSRecord(ctx context.Context, zoneID string, record models.DNSRecordCreateRequest) (*models.DNSRecord, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("zone ID is required")
	}

	log.Printf("Creating DNS record for zone: %s", zoneID)

	// Set the zone ID in the request
	record.ZoneID = cloudflare.F(zoneID)

	dnsRecord, err := s.client.DNS.Records.New(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS record: %w", err)
	}

	log.Printf("Successfully created DNS record: %s", dnsRecord.ID)
	return dnsRecord, nil
}

// CreateCNAMERecord creates a generic CNAME record for a hostname pointing to a target
func (s *CloudflareService) CreateCNAMERecord(ctx context.Context, zoneID, hostname, target string, proxied bool) (*models.DNSRecord, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("zone ID is required")
	}
	if hostname == "" {
		return nil, fmt.Errorf("hostname is required")
	}
	if target == "" {
		return nil, fmt.Errorf("target is required")
	}

	log.Printf("Creating CNAME record for hostname: %s -> %s (proxied: %t)", hostname, target, proxied)

	// Create CNAME record pointing to the target
	cname := dns.CNAMERecordParam{
		Type:    cloudflare.F(dns.CNAMERecordTypeCNAME),
		Name:    cloudflare.F(hostname),
		Content: cloudflare.F(target),
		TTL:     cloudflare.F(dns.TTL(1)), // Automatic TTL
		Proxied: cloudflare.F(proxied),    // Proxying through Cloudflare based on parameter
	}

	record := dns.RecordNewParams{
		ZoneID: cloudflare.F(zoneID),
		Body:   cname,
	}

	dnsRecord, err := s.client.DNS.Records.New(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to create CNAME record for hostname %s: %w", hostname, err)
	}

	log.Printf("Successfully created CNAME record: %s -> %s", hostname, target)
	return dnsRecord, nil
}

// GetDNSRecords retrieves DNS records for a zone
func (s *CloudflareService) GetDNSRecords(ctx context.Context, zoneID string) ([]models.DNSRecord, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("zone ID is required")
	}

	log.Printf("Fetching DNS records for zone: %s", zoneID)

	params := dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
	}

	records, err := s.client.DNS.Records.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DNS records: %w", err)
	}

	result := make([]models.DNSRecord, len(records.Result))
	for i, record := range records.Result {
		result[i] = record
	}

	log.Printf("Successfully fetched %d DNS records for zone %s", len(result), zoneID)
	return result, nil
}

// GetDNSRecordsByName retrieves DNS records by name for a zone
func (s *CloudflareService) GetDNSRecordsByName(ctx context.Context, zoneID, name string) ([]models.DNSRecord, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("zone ID is required")
	}
	if name == "" {
		return nil, fmt.Errorf("record name is required")
	}

	log.Printf("Fetching DNS records for zone: %s, name: %s", zoneID, name)

	params := dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
		Name:   cloudflare.F(dns.RecordListParamsName{Exact: cloudflare.F(name)}),
	}

	records, err := s.client.DNS.Records.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DNS records: %w", err)
	}

	result := make([]models.DNSRecord, len(records.Result))
	for i, record := range records.Result {
		result[i] = record
	}

	log.Printf("Successfully fetched %d DNS records for zone %s, name %s", len(result), zoneID, name)
	return result, nil
}

// UpdateDNSRecord updates an existing DNS record
func (s *CloudflareService) UpdateDNSRecord(ctx context.Context, zoneID, recordID string, record models.DNSRecordUpdateRequest) (*models.DNSRecord, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("zone ID is required")
	}
	if recordID == "" {
		return nil, fmt.Errorf("record ID is required")
	}

	log.Printf("Updating DNS record: %s in zone: %s", recordID, zoneID)

	// Set the zone ID in the request
	record.ZoneID = cloudflare.F(zoneID)

	dnsRecord, err := s.client.DNS.Records.Update(ctx, recordID, record)
	if err != nil {
		return nil, fmt.Errorf("failed to update DNS record: %w", err)
	}

	log.Printf("Successfully updated DNS record: %s", recordID)
	return dnsRecord, nil
}

// DeleteDNSRecord deletes a DNS record
func (s *CloudflareService) DeleteDNSRecord(ctx context.Context, zoneID, recordID string) error {
	if zoneID == "" {
		return fmt.Errorf("zone ID is required")
	}
	if recordID == "" {
		return fmt.Errorf("record ID is required")
	}

	log.Printf("Deleting DNS record: %s from zone: %s", recordID, zoneID)

	params := dns.RecordDeleteParams{
		ZoneID: cloudflare.F(zoneID),
	}

	_, err := s.client.DNS.Records.Delete(ctx, recordID, params)
	if err != nil {
		return fmt.Errorf("failed to delete DNS record: %w", err)
	}

	log.Printf("Successfully deleted DNS record: %s", recordID)
	return nil
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
