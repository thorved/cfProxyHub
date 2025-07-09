package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"cfPorxyHub/internal/models"
	"cfPorxyHub/internal/services"
	"cfPorxyHub/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CloudflareZoneHandler handles zone-related HTTP requests
type CloudflareZoneHandler struct {
	cfService *services.CloudflareService
}

// NewCloudflareZoneHandler creates a new zone handler instance
func NewCloudflareZoneHandler(cfService *services.CloudflareService) *CloudflareZoneHandler {
	return &CloudflareZoneHandler{
		cfService: cfService,
	}
}

// GetZonesByAccountID handles GET /api/cloudflare/accounts/{accountId}/zones
func (h *CloudflareZoneHandler) GetZonesByAccountID(c *gin.Context) {
	accountID := c.Param("accountId")
	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}

	// Get optional query parameters
	activeOnly := c.Query("active_only") == "true"
	searchTerm := c.Query("search")
	summaryOnly := c.Query("summary") == "true"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var zones []models.Zone
	var err error

	// Choose the appropriate service method based on query parameters
	if searchTerm != "" {
		zones, err = h.cfService.SearchZones(ctx, accountID, searchTerm)
	} else if activeOnly {
		zones, err = h.cfService.GetActiveZones(ctx, accountID)
	} else {
		zones, err = h.cfService.GetZonesByAccountID(ctx, accountID)
	}

	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch zones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return summary or full data based on query parameter
	if summaryOnly {
		summaries := models.NewZoneSummariesFromZones(zones)
		response := models.ZoneSummaryResponse{
			Success: true,
			Message: "Zones retrieved successfully",
			Data:    summaries,
			Total:   len(summaries),
		}
		utils.SuccessResponse(c, response)
	} else {
		response := models.ZonesResponse{
			Success: true,
			Message: "Zones retrieved successfully",
			Data:    zones,
			Total:   len(zones),
		}
		utils.SuccessResponse(c, response)
	}
}

// GetZoneByID handles GET /api/cloudflare/zones/{zoneId}
func (h *CloudflareZoneHandler) GetZoneByID(c *gin.Context) {
	zoneID := c.Param("zoneId")
	if zoneID == "" {
		utils.ErrorResponse(c, "Zone ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	zone, err := h.cfService.GetZoneByID(ctx, zoneID)
	if err != nil {
		utils.ErrorResponse(c, "Zone not found: "+err.Error(), http.StatusNotFound)
		return
	}

	response := models.ZoneResponse{
		Success: true,
		Message: "Zone retrieved successfully",
		Data:    *zone,
	}
	utils.SuccessResponse(c, response)
}

// GetZoneByName handles GET /api/cloudflare/accounts/{accountId}/zones/by-name/{domainName}
func (h *CloudflareZoneHandler) GetZoneByName(c *gin.Context) {
	accountID := c.Param("accountId")
	domainName := c.Param("domainName")

	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}
	if domainName == "" {
		utils.ErrorResponse(c, "Domain name is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	zone, err := h.cfService.GetZoneByName(ctx, accountID, domainName)
	if err != nil {
		utils.ErrorResponse(c, "Zone not found: "+err.Error(), http.StatusNotFound)
		return
	}

	response := models.ZoneResponse{
		Success: true,
		Message: "Zone retrieved successfully",
		Data:    *zone,
	}
	utils.SuccessResponse(c, response)
}

// GetZonesForDropdown handles GET /api/cloudflare/accounts/{accountId}/zones/dropdown
// Returns zone summaries optimized for dropdown usage
func (h *CloudflareZoneHandler) GetZonesForDropdown(c *gin.Context) {
	accountID := c.Param("accountId")
	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}

	// Get optional query parameters
	searchTerm := c.Query("search")
	limitStr := c.Query("limit")
	activeOnly := c.Query("active_only") != "false" // Default to true for dropdown

	// Parse limit parameter
	limit := 50 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var zones []models.Zone
	var err error

	// Choose the appropriate service method
	if searchTerm != "" {
		zones, err = h.cfService.SearchZones(ctx, accountID, searchTerm)
	} else if activeOnly {
		zones, err = h.cfService.GetActiveZones(ctx, accountID)
	} else {
		zones, err = h.cfService.GetZonesByAccountID(ctx, accountID)
	}

	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch zones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply limit
	if len(zones) > limit {
		zones = zones[:limit]
	}

	// Convert to summaries for dropdown
	summaries := models.NewZoneSummariesFromZones(zones)

	// Use standard response format for consistency
	utils.SuccessResponse(c, gin.H{
		"message": "Zones retrieved successfully",
		"zones":   summaries,
		"total":   len(summaries),
	})
}
