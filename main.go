package main

import (
	"context"
	"fmt"

	"github.com/tlop503/quicksand/docker_sdk"
)

func main() {
	ctx := context.Background()

	err := docker_sdk.StartContainer("jlesage/firefox", ctx, "firefox_go")
	if err != nil {
		panic(err)
	}

	err = docker_sdk.StopContainer(ctx, "containerName")
	if err != nil {
		fmt.Printf("Error stopping container: %v\n", err)
	} else {
		fmt.Println("Container stopped successfully")
	}
}
