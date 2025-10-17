package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/image"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

func main() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Pull the image if not present
	pullResp, err := cli.ImagePull(ctx, "jlesage/firefox", image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer pullResp.Close()

	// Drain the stream so the pull completes
	_, err = io.Copy(os.Stdout, pullResp)
	if err != nil {
		panic(err)
	}

	// Define port bindings
	portSet := nat.PortSet{
		"5800/tcp": struct{}{},
		"5900/tcp": struct{}{},
	}

	// map container ports to localhost ports
	portMap := nat.PortMap{
		"5800/tcp": []nat.PortBinding{{HostIP: "localhost", HostPort: "5800"}},
		"5900/tcp": []nat.PortBinding{{HostIP: "localhost", HostPort: "5900"}},
	}

	// Create the container
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:        "jlesage/firefox",
			ExposedPorts: portSet,
		},
		&container.HostConfig{
			PortBindings: portMap,
		},
		&network.NetworkingConfig{},
		nil,
		"firefox-go",
	)
	if err != nil {
		panic(err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("Container started: ", resp.ID)

	// Attach to logs
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Print logs to stdout
	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		panic(err)
	}
}
