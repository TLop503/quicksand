package main

import (
	"context"
	"fmt"

	"github.com/tlop503/quicksand/docker_sdk"
)

func main() {
	ctx := context.Background()

	containerName, err := docker_sdk.StartContainer("jlesage/firefox", ctx, "firefox_go")
	if err != nil {
		panic(err)
	}

	err = docker_sdk.StopContainer(ctx, containerName)
	if err != nil {
		fmt.Printf("Error stopping container: %v\n", err)
	} else {
		fmt.Println("Container stopped successfully")
	}

	// Remove the container (using force=true in case stop failed)
	err = docker_sdk.RemoveContainer(ctx, containerName, false)
	if err != nil {
		fmt.Printf("Error removing container: %v\n", err)
	} else {
		fmt.Println("Container removed successfully")
	}
}
