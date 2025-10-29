package docker_sdk

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

func StartContainer(img string, ctx context.Context, ctrName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", errors.Errorf("Failed to create client: %v", err)
	}

	// pull image if not present
	// img is name of container (e.x. jlesage/firefox)
	// image.PullOptions are vars for how pulling occurs
	pullResponse, err := cli.ImagePull(ctx, img, types.ImagePullOptions{})
	if err != nil {
		return "", errors.Errorf("Failed to pull image: %v", err)
	}
	defer pullResponse.Close()

	// flush stream to complete pull
	// TODO: confirm if this is needed
	_, err = io.Copy(os.Stdout, pullResponse)
	if err != nil {
		return "", errors.Errorf("Failed to pull image and flush io: %v", err)
	}

	// port bindings!
	portSet := nat.PortSet{
		"5800/tcp": struct{}{},
		"5900/tcp": struct{}{},
	}
	portMap := nat.PortMap{
		"5800/tcp": []nat.PortBinding{{HostIP: "localhost", HostPort: "5800"}},
		"5900/tcp": []nat.PortBinding{{HostIP: "localhost", HostPort: "5900"}},
	}

	// Generate the randomized name and store it
	actualName := randomizeName(ctrName)

	// Create container
	containerResponse, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:        img,
			ExposedPorts: portSet,
		},
		&container.HostConfig{
			PortBindings: portMap,
		},
		&network.NetworkingConfig{},
		nil,
		actualName,
	)
	if err != nil {
		return "", errors.Errorf("Failed to create container: %v", err)
	}

	// Start the container!
	if err := cli.ContainerStart(ctx, containerResponse.ID, container.StartOptions{}); err != nil {
		return "", errors.Errorf("Failed to start container: %v", err)
	}

	fmt.Printf("Started container %s", img)

	return actualName, nil
}

// randomize names for containers to create unique options
func randomizeName(ctrName string) string {
	timestamp := time.Now().Unix()                // seconds since epoch
	last4 := fmt.Sprintf("%04d", timestamp%10000) // get last 4 digits
	return ctrName + last4
}

func StopContainer(ctx context.Context, ctrName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Errorf("Failed to create Docker client: %v", err)
	}

	// Attempt to stop the container, using the standard stop time
	err = cli.ContainerStop(ctx, ctrName, container.StopOptions{})
	if err != nil {
		return errors.Errorf("Failed to stop container %s: %v", ctrName, err)
	}

	fmt.Printf("Stopped container %s\n", ctrName)
	return nil
}

func RemoveContainer(ctx context.Context, ctrName string, force bool) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Errorf("Failed to create Docker client: %v", err)
	}

	// Check if container exists and is stopped before removing
	// Should add a call to stop container in this if container running?
	if !force {
		containerInfo, err := cli.ContainerInspect(ctx, ctrName)
		if err != nil {
			return errors.Errorf("Container %s not found: %v", ctrName, err)
		}

		if containerInfo.State.Running {
			return errors.Errorf("Container %s is still running. Stop it first or use force=true", ctrName)
		}
	}

	// Struct that manages the type to removal
	removeOptions := types.ContainerRemoveOptions{
		Force:         force,
		RemoveVolumes: true,
		RemoveLinks:   false, // explicitly set to false unless you want link removal
	}

	// Remove container
	err = cli.ContainerRemove(ctx, ctrName, removeOptions)
	if err != nil {
		return errors.Errorf("Failed to remove container %s: %v", ctrName, err)
	}

	fmt.Printf("Removed container %s\n", ctrName)
	return nil
}
