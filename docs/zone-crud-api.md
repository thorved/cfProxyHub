# Zone CRUD API Endpoints

This document describes the CRUD (Create, Read, Update, Delete) operations available for Cloudflare zones.

## Authentication
All endpoints require authentication via the `RequireAuth()` middleware.

## Base URL
All endpoints are prefixed with `/api/cloudflare`

## Endpoints

### 1. Create Zone
**POST** `/accounts/{accountId}/zones`

Creates a new zone (domain) in the specified Cloudflare account.

#### Request Body
```json
{
  "name": "example.com",
  "account_id": "your-account-id"
}
```

#### Response
```json
{
  "success": true,
  "message": "Zone created successfully",
  "data": {
    "id": "zone-id",
    "name": "example.com",
    "status": "pending",
    "type": "full",
    "development_mode": 0,
    "name_servers": ["ns1.cloudflare.com", "ns2.cloudflare.com"],
    "original_name_servers": [],
    "original_registrar": "",
    "original_dnshost": "",
    "modified_on": "2025-07-09T10:00:00Z",
    "created_on": "2025-07-09T10:00:00Z",
    "activated_on": "0001-01-01T00:00:00Z"
  }
}
```

### 2. Read Zone by ID
**GET** `/zones/{zoneId}`

Retrieves a specific zone by its ID.

#### Response
```json
{
  "success": true,
  "message": "Zone retrieved successfully",
  "data": {
    "id": "zone-id",
    "name": "example.com",
    "status": "active",
    // ... other zone properties
  }
}
```

### 3. Read Zones by Account
**GET** `/accounts/{accountId}/zones`

Retrieves all zones for a specific account with optional filtering.

#### Query Parameters
- `active_only` (boolean): Filter to only active zones
- `search` (string): Search term to filter zones by name
- `summary` (boolean): Return only summary information

#### Response
```json
{
  "success": true,
  "message": "Zones retrieved successfully",
  "data": [
    {
      "id": "zone-id-1",
      "name": "example1.com",
      "status": "active",
      // ... other zone properties
    },
    {
      "id": "zone-id-2",
      "name": "example2.com",
      "status": "active",
      // ... other zone properties
    }
  ],
  "total": 2
}
```

### 4. Update Zone
**PUT** `/zones/{zoneId}`

Updates an existing zone. Currently supports pausing/unpausing zones.

#### Request Body
```json
{
  "paused": false
}
```

#### Response
```json
{
  "success": true,
  "message": "Zone updated successfully",
  "data": {
    "id": "zone-id",
    "name": "example.com",
    "status": "active",
    // ... other zone properties
  }
}
```

### 5. Delete Zone
**DELETE** `/zones/{zoneId}`

Deletes a zone permanently.

#### Response
```json
{
  "success": true,
  "message": "Zone deleted successfully",
  "id": "zone-id"
}
```

### 6. Additional Read Operations

#### Get Zone by Name
**GET** `/accounts/{accountId}/zones/by-name/{domainName}`

Retrieves a zone by its domain name within a specific account.

#### Get Zones for Dropdown
**GET** `/accounts/{accountId}/zones/dropdown`

Returns zone summaries optimized for dropdown usage.

#### Query Parameters
- `search` (string): Search term to filter zones
- `limit` (number): Maximum number of results (default: 50, max: 100)
- `active_only` (boolean): Filter to only active zones (default: true)

## Error Responses

All endpoints return error responses in the following format when something goes wrong:

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error information"
}
```

Common HTTP status codes:
- `400 Bad Request`: Invalid request parameters or body
- `401 Unauthorized`: Authentication required
- `404 Not Found`: Zone not found
- `500 Internal Server Error`: Server-side error

## Usage Examples

### Creating a new zone
```bash
curl -X POST "http://localhost:8080/api/cloudflare/accounts/your-account-id/zones" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "name": "mynewdomain.com",
    "account_id": "your-account-id"
  }'
```

### Pausing a zone
```bash
curl -X PUT "http://localhost:8080/api/cloudflare/zones/zone-id" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "paused": true
  }'
```

### Deleting a zone
```bash
curl -X DELETE "http://localhost:8080/api/cloudflare/zones/zone-id" \
  -H "Authorization: Bearer your-token"
```
