package models

import (
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// Tunnel represents a Cloudflare tunnel with additional metadata
type Tunnel struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Secret         string             `json:"secret,omitempty"`
	CreatedAt      *time.Time         `json:"created_at,omitempty"`
	DeletedAt      *time.Time         `json:"deleted_at,omitempty"`
	Connections    []TunnelConnection `json:"connections,omitempty"`
	ConnsActiveAt  *time.Time         `json:"conns_active_at,omitempty"`
	ConnInactiveAt *time.Time         `json:"conn_inactive_at,omitempty"`
	TunnelType     string             `json:"tunnel_type,omitempty"`
	Status         string             `json:"status,omitempty"`
	RemoteConfig   bool               `json:"remote_config,omitempty"`
}

// TunnelConnection represents a connection to a tunnel
type TunnelConnection struct {
	ColoName           string `json:"colo_name"`
	ID                 string `json:"id"`
	IsPendingReconnect bool   `json:"is_pending_reconnect"`
	ClientID           string `json:"client_id"`
	ClientVersion      string `json:"client_version"`
	OpenedAt           string `json:"opened_at"`
	OriginIP           string `json:"origin_ip"`
}

// TunnelListParams represents parameters for listing tunnels
type TunnelListParams struct {
	Name          string     `json:"name,omitempty"`
	IsDeleted     *bool      `json:"is_deleted,omitempty"`
	ExcludePrefix string     `json:"exclude_prefix,omitempty"`
	ExistedAt     *time.Time `json:"existed_at,omitempty"`
	IncludePrefix string     `json:"include_prefix,omitempty"`
	UUID          string     `json:"uuid,omitempty"`
}

// TunnelCreateRequest represents the request structure for creating a tunnel
type TunnelCreateRequest struct {
	Name       string `json:"name" binding:"required"`
	Secret     string `json:"secret,omitempty"`
	ConfigSrc  string `json:"config_src,omitempty"`
	TunnelType string `json:"tunnel_type,omitempty"`
}

// TunnelUpdateRequest represents the request structure for updating a tunnel
type TunnelUpdateRequest struct {
	Name       string `json:"name,omitempty"`
	Secret     string `json:"secret,omitempty"`
	ConfigSrc  string `json:"config_src,omitempty"`
	TunnelType string `json:"tunnel_type,omitempty"`
}

// TunnelResponse represents the response structure for single tunnel endpoint
type TunnelResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Tunnel `json:"data"`
}

// TunnelsResponse represents the response structure for tunnels list endpoint
type TunnelsResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []Tunnel `json:"data"`
	Total   int      `json:"total"`
}

// TunnelToken represents a tunnel token
type TunnelToken struct {
	Token string `json:"token"`
}

// TunnelTokenResponse represents the response structure for tunnel token endpoint
type TunnelTokenResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    TunnelToken `json:"data"`
}

// ConvertFromCloudflareTunnel converts a cloudflare.Tunnel to our Tunnel model
func ConvertFromCloudflareTunnel(cfTunnel cloudflare.Tunnel) Tunnel {
	tunnel := Tunnel{
		ID:             cfTunnel.ID,
		Name:           cfTunnel.Name,
		Secret:         cfTunnel.Secret,
		CreatedAt:      cfTunnel.CreatedAt,
		DeletedAt:      cfTunnel.DeletedAt,
		ConnsActiveAt:  cfTunnel.ConnsActiveAt,
		ConnInactiveAt: cfTunnel.ConnInactiveAt,
		TunnelType:     cfTunnel.TunnelType,
		Status:         cfTunnel.Status,
		RemoteConfig:   cfTunnel.RemoteConfig,
	}

	// Convert connections
	if cfTunnel.Connections != nil {
		tunnel.Connections = make([]TunnelConnection, len(cfTunnel.Connections))
		for i, conn := range cfTunnel.Connections {
			tunnel.Connections[i] = TunnelConnection{
				ColoName:           conn.ColoName,
				ID:                 conn.ID,
				IsPendingReconnect: conn.IsPendingReconnect,
				ClientID:           conn.ClientID,
				ClientVersion:      conn.ClientVersion,
				OpenedAt:           conn.OpenedAt,
				OriginIP:           conn.OriginIP,
			}
		}
	}

	return tunnel
}

// ConvertToCloudflareTunnel converts our Tunnel model to cloudflare.Tunnel
func (t *Tunnel) ConvertToCloudflareTunnel() cloudflare.Tunnel {
	cfTunnel := cloudflare.Tunnel{
		ID:             t.ID,
		Name:           t.Name,
		Secret:         t.Secret,
		CreatedAt:      t.CreatedAt,
		DeletedAt:      t.DeletedAt,
		ConnsActiveAt:  t.ConnsActiveAt,
		ConnInactiveAt: t.ConnInactiveAt,
		TunnelType:     t.TunnelType,
		Status:         t.Status,
		RemoteConfig:   t.RemoteConfig,
	}

	// Convert connections
	if t.Connections != nil {
		cfTunnel.Connections = make([]cloudflare.TunnelConnection, len(t.Connections))
		for i, conn := range t.Connections {
			cfTunnel.Connections[i] = cloudflare.TunnelConnection{
				ColoName:           conn.ColoName,
				ID:                 conn.ID,
				IsPendingReconnect: conn.IsPendingReconnect,
				ClientID:           conn.ClientID,
				ClientVersion:      conn.ClientVersion,
				OpenedAt:           conn.OpenedAt,
				OriginIP:           conn.OriginIP,
			}
		}
	}

	return cfTunnel
}

// ConvertFromCloudflareTunnelListParams converts cloudflare.TunnelListParams to our TunnelListParams
func ConvertFromCloudflareTunnelListParams(cfParams cloudflare.TunnelListParams) TunnelListParams {
	return TunnelListParams{
		Name:          cfParams.Name,
		IsDeleted:     cfParams.IsDeleted,
		ExcludePrefix: cfParams.ExcludePrefix,
		ExistedAt:     cfParams.ExistedAt,
		IncludePrefix: cfParams.IncludePrefix,
		UUID:          cfParams.UUID,
	}
}

// ConvertToCloudflareTunnelListParams converts our TunnelListParams to cloudflare.TunnelListParams
func (p *TunnelListParams) ConvertToCloudflareTunnelListParams() cloudflare.TunnelListParams {
	return cloudflare.TunnelListParams{
		Name:          p.Name,
		IsDeleted:     p.IsDeleted,
		ExcludePrefix: p.ExcludePrefix,
		ExistedAt:     p.ExistedAt,
		IncludePrefix: p.IncludePrefix,
		UUID:          p.UUID,
	}
}
