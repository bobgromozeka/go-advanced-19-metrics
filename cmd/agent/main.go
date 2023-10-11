package main

import (
	"github.com/bobgromozeka/metrics/internal/agent"
)

func main() {
	printMetadata()
	setupConfiguration()

	agent.Run(startupConfig)
}
