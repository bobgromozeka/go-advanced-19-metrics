package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var serverAddr string
var pollInterval int
var reportInterval int

const (
	Address        = "ADDRESS"
	ReportInterval = "REPORT_INTERVAL"
	PollInterval   = "POLL_INTERVAL"
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "server address to send metrics")
	flag.IntVar(&pollInterval, "p", 2, "Metrics polling interval")
	flag.IntVar(&reportInterval, "r", 10, "Metrics reporting interval to server")

	flag.Parse()
}

func parseEnv() {
	if addr := os.Getenv(Address); addr != "" {
		serverAddr = addr
	}

	if ri := os.Getenv(ReportInterval); ri != "" {
		parsedRi, err := strconv.Atoi(ri)
		if err != nil {
			fmt.Println("You've specified wrong report interval. Should be integer, got: ", ri)
		}

		reportInterval = parsedRi
	}

	if pi := os.Getenv(PollInterval); pi != "" {
		parsedPi, err := strconv.Atoi(pi)
		if err != nil {
			fmt.Println("You've specified wrong poll interval. Should be integer, got: ", pi)
		}

		pollInterval = parsedPi
	}
}

func setupConfiguration() {
	parseFlags()
	parseEnv()
}
