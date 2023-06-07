package main

import "github.com/bobgromozeka/metrics/internal/agent"

func main() {
	setupConfiguration()

	agent.Run(agent.StartupConfig{
		ServerAddr:     serverAddr,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	})
}
