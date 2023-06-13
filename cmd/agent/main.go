package main

import "github.com/bobgromozeka/metrics/internal/agent"

func main() {
	setupConfiguration()

	agent.Run(agent.StartupConfig{
		ServerAddr:     serverAddr,
		ServerScheme:   serverScheme,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	})
}
