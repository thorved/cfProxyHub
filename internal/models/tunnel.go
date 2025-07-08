package models

import (
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"
)

// Use v4 types directly - no conversion needed
type Tunnel = zero_trust.TunnelCloudflaredGetResponse
type TunnelListResponse = zero_trust.TunnelCloudflaredListResponse
type TunnelCreateRequest = zero_trust.TunnelCloudflaredNewParams
type TunnelUpdateRequest = zero_trust.TunnelCloudflaredEditParams
type TunnelListParams = zero_trust.TunnelCloudflaredListParams
type TunnelNewResponse = zero_trust.TunnelCloudflaredNewResponse
type TunnelEditResponse = zero_trust.TunnelCloudflaredEditResponse
type TunnelDeleteResponse = zero_trust.TunnelCloudflaredDeleteResponse

// Response wrappers for API endpoints
type TunnelResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Tunnel `json:"data"`
}

type TunnelsResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    []TunnelListResponse `json:"data"`
	Total   int                  `json:"total"`
}

type TunnelToken struct {
	Token string `json:"token"`
}

type TunnelTokenResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    TunnelToken `json:"data"`
}

// Helper functions for creating request parameters
func NewTunnelCreateRequest(name, configSrc string) TunnelCreateRequest {
	params := zero_trust.TunnelCloudflaredNewParams{
		Name: cloudflare.F(name),
	}
	if configSrc != "" {
		params.ConfigSrc = cloudflare.F(zero_trust.TunnelCloudflaredNewParamsConfigSrc(configSrc))
	}
	return params
}

func NewTunnelUpdateRequest(name string) TunnelUpdateRequest {
	params := zero_trust.TunnelCloudflaredEditParams{}
	if name != "" {
		params.Name = cloudflare.F(name)
	}
	return params
}

func NewTunnelListParams(name string) TunnelListParams {
	params := zero_trust.TunnelCloudflaredListParams{}
	if name != "" {
		params.Name = cloudflare.F(name)
	}
	return params
}
