package models

import (
	"time"

	"github.com/cloudflare/cloudflare-go/v4/dns"
)

// Zone represents a Cloudflare zone/domain
type Zone struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Status              string                 `json:"status"`
	Type                string                 `json:"type"`
	DevelopmentMode     int                    `json:"development_mode"`
	NameServers         []string               `json:"name_servers"`
	OriginalNameServers []string               `json:"original_name_servers"`
	OriginalRegistrar   string                 `json:"original_registrar"`
	OriginalDNSHost     string                 `json:"original_dnshost"`
	ModifiedOn          time.Time              `json:"modified_on"`
	CreatedOn           time.Time              `json:"created_on"`
	ActivatedOn         time.Time              `json:"activated_on"`
	Meta                map[string]interface{} `json:"meta"`
	Owner               map[string]interface{} `json:"owner"`
	Account             map[string]interface{} `json:"account"`
	Permissions         []string               `json:"permissions"`
	Plan                map[string]interface{} `json:"plan"`
}

// DNS Record types - use v4 types directly
type DNSRecord = dns.RecordResponse
type DNSRecordCreateRequest = dns.RecordNewParams
type DNSRecordUpdateRequest = dns.RecordUpdateParams
type DNSRecordDeleteResponse = dns.RecordDeleteResponse

// Response wrappers for API endpoints
type ZoneResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Zone   `json:"data"`
}

type ZonesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []Zone `json:"data"`
	Total   int    `json:"total"`
}

// DNS Record response wrappers
type DNSRecordResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    DNSRecord `json:"data"`
}

type DNSRecordsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    []DNSRecord `json:"data"`
	Total   int         `json:"total"`
}

// Zone summary for dropdowns
type ZoneSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

type ZoneSummaryResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []ZoneSummary `json:"data"`
	Total   int           `json:"total"`
}

// Zone create/update request models
type ZoneCreateRequest struct {
	Name      string `json:"name" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
}

type ZoneUpdateRequest struct {
	Paused bool `json:"paused"`
}

// Zone delete response
type ZoneDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

// Helper functions for creating zone summaries
func NewZoneSummary(zone Zone) ZoneSummary {
	return ZoneSummary{
		ID:     zone.ID,
		Name:   zone.Name,
		Status: zone.Status,
		Type:   zone.Type,
	}
}

func NewZoneSummariesFromZones(zoneList []Zone) []ZoneSummary {
	summaries := make([]ZoneSummary, len(zoneList))
	for i, zone := range zoneList {
		summaries[i] = NewZoneSummary(zone)
	}
	return summaries
}
