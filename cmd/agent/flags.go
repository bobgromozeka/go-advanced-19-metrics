package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/bobgromozeka/metrics/internal/agent"
)

var startupConfig agent.StartupConfig

const JSONConfigPath = "CONFIG"

const (
	Address        = "ADDRESS"
	ReportInterval = "REPORT_INTERVAL"
	PollInterval   = "POLL_INTERVAL"
	Key            = "KEY"
	PublicKeyPath  = "CRYPTO_KEY"
)

func parseFlags() {
	flag.StringVar(&startupConfig.ServerAddr, "a", "localhost:8080", "server address to send metrics")
	flag.StringVar(&startupConfig.ServerScheme, "s", "http", "server scheme (http, https)")
	flag.IntVar(&startupConfig.PollInterval, "p", 2, "Metrics polling interval")
	flag.IntVar(&startupConfig.ReportInterval, "r", 10, "Metrics reporting interval to server")
	flag.StringVar(&startupConfig.HashKey, "k", "", "Key to make request signature")
	flag.StringVar(&startupConfig.PublicKeyPath, "ck", "./public.pem", "Public key for data encryption")

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

	if key := os.Getenv(Key); key != "" {
		startupConfig.HashKey = key
	}

	if publicKeyPath := os.Getenv(PublicKeyPath); publicKeyPath != "" {
		startupConfig.PublicKeyPath = publicKeyPath
	}
}

func parseJSONConfig() {
	if os.Getenv(JSONConfigPath) == "" {
		return
	}

	conf, err := os.Open(JSONConfigPath)
	if err != nil {
		fmt.Printf("Could not open json config: %v \n", err)
	}

	decoder := json.NewDecoder(conf)

	if decodeErr := decoder.Decode(&startupConfig); decodeErr != nil {
		fmt.Printf("Could not open parse json config: %v \n", decodeErr)
	}
}

func setupConfiguration() {
	parseJSONConfig()
	parseFlags()
	parseEnv()
}
