package main

import (
	"context"

	"github.com/bobgromozeka/metrics/internal/agent"
	"github.com/bobgromozeka/metrics/internal/helpers"
)

func main() {
	printMetadata()
	setupConfiguration()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	helpers.SetupGracefulShutdown(cancel)

	agent.Run(ctx, startupConfig)
}
