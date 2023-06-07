package container

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient interface {
	NewClient() error
	ContainerList(ctx context.Context, opts types.ContainerListOptions) ([]types.Container, error)
	ImagePull(ctx context.Context, image string, opts types.ImagePullOptions) (io.ReadCloser, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, containerName string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	// ...
}

type DefaultDockerClient struct {
	*client.Client
}

func (dc *DefaultDockerClient) NewClient() error {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	dc.Client = client
	return nil
}

func (dc *DefaultDockerClient) ContainerList(ctx context.Context, opts types.ContainerListOptions) ([]types.Container, error) {
	return dc.Client.ContainerList(ctx, opts)
}

func (dc *DefaultDockerClient) ImagePull(ctx context.Context, image string, opts types.ImagePullOptions) (io.ReadCloser, error) {
	return dc.Client.ImagePull(ctx, image, opts)
}

func (dc *DefaultDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, containerName string) (container.CreateResponse, error) {
	return dc.Client.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
}

func (dc *DefaultDockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return dc.Client.ContainerStart(ctx, containerID, options)
}
