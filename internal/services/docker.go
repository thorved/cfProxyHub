package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// DockerService provides methods to interact with Docker
// You can choose to connect via API or socket by setting the Docker host
// Example: export DOCKER_HOST=unix:///var/run/docker.sock or DOCKER_HOST=tcp://localhost:2375
// See: https://docs.docker.com/engine/reference/commandline/cli/#environment-variables
type DockerService struct {
	cli *client.Client
}

// NewDockerService creates a new DockerService instance
func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &DockerService{cli: cli}, nil
}

// Ping checks if the Docker daemon is accessible
func (ds *DockerService) Ping(ctx context.Context) error {
	if ds.cli == nil {
		return fmt.Errorf("docker client not initialized")
	}

	_, err := ds.cli.Ping(ctx)
	return err
}

// ListContainers returns running containers
func (ds *DockerService) ListContainers(ctx context.Context) ([]types.Container, error) {
	return ds.cli.ContainerList(ctx, container.ListOptions{All: true})
}

// ListImages returns Docker images
func (ds *DockerService) ListImages(ctx context.Context) ([]image.Summary, error) {
	return ds.cli.ImageList(ctx, image.ListOptions{All: true})
}

// ListVolumes returns Docker volumes
func (ds *DockerService) ListVolumes(ctx context.Context) ([]*volume.Volume, error) {
	resp, err := ds.cli.VolumeList(ctx, volume.ListOptions{Filters: filters.Args{}})
	if err != nil {
		return nil, err
	}
	return resp.Volumes, nil
}

// ListNetworks returns Docker networks
func (ds *DockerService) ListNetworks(ctx context.Context) ([]network.Summary, error) {
	resp, err := ds.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateContainerParams contains the parameters for creating a Docker container
type CreateContainerParams struct {
	Name          string            `json:"name"`
	Image         string            `json:"image"`
	Ports         map[string]string `json:"ports"`   // host_port:container_port
	Volumes       map[string]string `json:"volumes"` // host_path:container_path
	Env           []string          `json:"env"`
	Command       []string          `json:"command,omitempty"`
	NetworkMode   string            `json:"network_mode,omitempty"`
	RestartPolicy string            `json:"restart_policy,omitempty"`
}

// CreateContainer creates a new container
func (ds *DockerService) CreateContainer(ctx context.Context, params CreateContainerParams) (string, error) {
	// Configure port bindings
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}

	for hostPort, containerPort := range params.Ports {
		port := nat.Port(containerPort)
		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			},
		}
	}

	// Configure volume bindings
	var binds []string
	for hostPath, containerPath := range params.Volumes {
		binds = append(binds, fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	// Configure restart policy
	restartPolicy := container.RestartPolicy{}
	if params.RestartPolicy != "" {
		restartPolicy = container.RestartPolicy{Name: container.RestartPolicyMode(params.RestartPolicy)}
	}

	// Create container
	resp, err := ds.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        params.Image,
			Cmd:          params.Command,
			Env:          params.Env,
			ExposedPorts: exposedPorts,
		},
		&container.HostConfig{
			PortBindings:  portBindings,
			Binds:         binds,
			NetworkMode:   container.NetworkMode(params.NetworkMode),
			RestartPolicy: restartPolicy,
		},
		nil,
		nil,
		params.Name,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return resp.ID, nil
}

// RemoveContainer removes a container by ID
func (ds *DockerService) RemoveContainer(ctx context.Context, id string) error {
	return ds.cli.ContainerRemove(ctx, id, container.RemoveOptions{
		Force: true,
	})
}

// StartContainer starts a container by ID
func (ds *DockerService) StartContainer(ctx context.Context, id string) error {
	return ds.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// StopContainer stops a container by ID
func (ds *DockerService) StopContainer(ctx context.Context, id string) error {
	timeout := 30 // seconds
	return ds.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
}

// CloudflareTunnelParams contains parameters for creating a Cloudflare Tunnel container
type CloudflareTunnelParams struct {
	Name          string `json:"name"`
	Token         string `json:"token"`
	RestartPolicy string `json:"restart_policy,omitempty"`
}

// CreateCloudflareTunnelContainer creates a container running a Cloudflare Tunnel
func (ds *DockerService) CreateCloudflareTunnelContainer(ctx context.Context, params CloudflareTunnelParams) (string, error) {
	// Set default container name if not provided
	containerName := params.Name
	if containerName == "" {
		containerName = "cloudflared-tunnel"
	}

	// Ensure the name has the right prefix if not already present
	if !strings.HasPrefix(containerName, "cloudflared-") &&
		!strings.HasPrefix(containerName, "cloudflare-") {
		containerName = "cloudflared-" + containerName
	}

	// Configure restart policy
	restartPolicy := container.RestartPolicy{}
	if params.RestartPolicy != "" {
		restartPolicy = container.RestartPolicy{Name: container.RestartPolicyMode(params.RestartPolicy)}
	} else {
		// Default to always restart
		restartPolicy = container.RestartPolicy{Name: container.RestartPolicyAlways}
	}

	// Command to run the tunnel with the provided token
	cmd := []string{"tunnel", "--no-autoupdate", "run", "--token", params.Token}

	// Pull the image first to ensure we have the latest
	_, err := ds.cli.ImagePull(ctx, "cloudflare/cloudflared:latest", image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull Cloudflare cloudflared image: %w", err)
	}

	// Create the Cloudflare tunnel container with multiple labels
	resp, err := ds.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "cloudflare/cloudflared:latest",
			Cmd:   cmd,
			Labels: map[string]string{
				"com.cloudflare.tunnel": "true",
				"app":                   "cloudflared",
				"service":               "tunnel",
				"managed-by":            "cfproxyhub",
			},
		},
		&container.HostConfig{
			RestartPolicy: restartPolicy,
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create Cloudflare tunnel container: %w", err)
	}

	return resp.ID, nil
}

// FindCloudflareTunnelContainers finds all containers running Cloudflare tunnels
func (ds *DockerService) FindCloudflareTunnelContainers(ctx context.Context) ([]types.Container, error) {
	// Try first with the label filter
	args := filters.NewArgs()
	args.Add("label", "com.cloudflare.tunnel=true")

	// List containers with the filter
	labeledContainers, err := ds.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: args,
	})

	if err != nil {
		return nil, fmt.Errorf("error listing labeled containers: %w", err)
	}

	// If we found labeled containers, return them
	if len(labeledContainers) > 0 {
		return labeledContainers, nil
	}

	// As a fallback, look for cloudflare container names or images
	allContainers, err := ds.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})

	if err != nil {
		return nil, fmt.Errorf("error listing all containers: %w", err)
	}

	var tunnelContainers []types.Container

	// Filter containers by image or name
	for _, c := range allContainers {
		// Check if image contains cloudflare or cloudflared
		if c.Image == "cloudflare/cloudflared:latest" ||
			c.Image == "cloudflare/cloudflared" {
			tunnelContainers = append(tunnelContainers, c)
			continue
		}

		// Check container names
		for _, name := range c.Names {
			if name != "" && (contains(name, "cloudflare") ||
				contains(name, "cloudflared") ||
				contains(name, "tunnel")) {
				tunnelContainers = append(tunnelContainers, c)
				break
			}
		}
	}

	return tunnelContainers, nil
}

// InspectContainer gets detailed information about a container
func (ds *DockerService) InspectContainer(ctx context.Context, id string) (types.ContainerJSON, error) {
	return ds.cli.ContainerInspect(ctx, id)
}

// helper function to check if a string contains another string (case insensitive)
func contains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
