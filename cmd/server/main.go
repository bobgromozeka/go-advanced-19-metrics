package main

import (
	"context"
	"fmt"

	"github.com/bobgromozeka/metrics/internal/helpers"
	"github.com/bobgromozeka/metrics/internal/server"
)

func main() {
	printMetadata()
	setupConfiguration()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	helpers.SetupGracefulShutdown(cancel)

	err := server.Start(ctx, startupConfig)

	if err != nil {
		fmt.Println("Error during server start: ", err)
	}
}
