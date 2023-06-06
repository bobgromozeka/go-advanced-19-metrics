package main

import "github.com/bobgromozeka/metrics/internal/agent"

func main() {
	parseFlags()

	agent.Run(agent.StartupConfig{
		ServerAddr:     serverAddr,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	})
}
