package main

import (
	"context"
	"github.com/tlop503/quicksand/docker_sdk"
)

func main() {
	ctx := context.Background()

	err := docker_sdk.StartContainer("jlesage/firefox", ctx, "firefox_go")
	if err != nil {
		panic(err)
	}
}
