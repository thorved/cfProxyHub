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

// Public Hostname types - these are managed through tunnel configuration
type TunnelConfiguration = zero_trust.TunnelCloudflaredConfigurationGetResponse
type TunnelConfigurationUpdate = zero_trust.TunnelCloudflaredConfigurationUpdateParams
type TunnelConfigurationUpdateResponse = zero_trust.TunnelCloudflaredConfigurationUpdateResponse
type PublicHostnameIngress = zero_trust.TunnelCloudflaredConfigurationGetResponseConfigIngress
type PublicHostnameIngressParam = zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress

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

// Public Hostname response wrappers
type PublicHostnameResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    PublicHostnameIngress `json:"data"`
}

type PublicHostnamesResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    []PublicHostnameIngress `json:"data"`
	Total   int                     `json:"total"`
}

type TunnelConfigurationResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    TunnelConfiguration `json:"data"`
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

// Helper functions for creating public hostname parameters
func NewPublicHostnameIngressParam(hostname, service, path string) PublicHostnameIngressParam {
	param := zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress{
		Hostname: cloudflare.F(hostname),
		Service:  cloudflare.F(service),
	}
	if path != "" {
		param.Path = cloudflare.F(path)
	}
	return param
}

func NewPublicHostnameIngressParamWithOriginRequest(hostname, service, path string, originRequest zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngressOriginRequest) PublicHostnameIngressParam {
	param := zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress{
		Hostname:      cloudflare.F(hostname),
		Service:       cloudflare.F(service),
		OriginRequest: cloudflare.F(originRequest),
	}
	if path != "" {
		param.Path = cloudflare.F(path)
	}
	return param
}
