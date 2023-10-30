package main

import (
	"context"

	"github.com/bobgromozeka/metrics/internal/helpers"
	"github.com/bobgromozeka/metrics/internal/server"
)

func main() {
	printMetadata()
	setupConfiguration()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	helpers.SetupGracefulShutdown(cancel)

	server.Start(ctx, startupConfig)
}
