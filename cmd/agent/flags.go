package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/bobgromozeka/metrics/internal/agent"
)

var startupConfig agent.StartupConfig

const (
	Address        = "ADDRESS"
	ReportInterval = "REPORT_INTERVAL"
	PollInterval   = "POLL_INTERVAL"
)

func parseFlags() {
	flag.StringVar(&startupConfig.ServerAddr, "a", "localhost:8080", "server address to send metrics")
	flag.StringVar(&startupConfig.ServerScheme, "s", "http", "server scheme (http, https)")
	flag.IntVar(&startupConfig.PollInterval, "p", 2, "Metrics polling interval")
	flag.IntVar(&startupConfig.ReportInterval, "r", 10, "Metrics reporting interval to server")

	flag.Parse()
}

func parseEnv() {
	if addr := os.Getenv(Address); addr != "" {
		startupConfig.ServerAddr = addr
	}

	if ri := os.Getenv(ReportInterval); ri != "" {
		parsedRi, err := strconv.Atoi(ri)
		if err != nil {
			fmt.Println("You've specified wrong report interval. Should be integer, got: ", ri)
		}

		startupConfig.ReportInterval = parsedRi
	}

	if pi := os.Getenv(PollInterval); pi != "" {
		parsedPi, err := strconv.Atoi(pi)
		if err != nil {
			fmt.Println("You've specified wrong poll interval. Should be integer, got: ", pi)
		}

		startupConfig.PollInterval = parsedPi
	}
}

func setupConfiguration() {
	parseFlags()
	parseEnv()
}
