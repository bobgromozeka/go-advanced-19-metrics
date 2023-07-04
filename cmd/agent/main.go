package main

import (
	"github.com/bobgromozeka/metrics/internal/agent"
)

func main() {
	setupConfiguration()

	agent.Run(startupConfig)
}
