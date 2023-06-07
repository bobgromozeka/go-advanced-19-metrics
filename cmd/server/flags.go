package main

import (
	"flag"
	"os"
)

var serverAddr string

const (
	Address = "ADDRESS"
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", ":8080", "address and port to run server")

	flag.Parse()
}

func parseEnv() {
	if addr := os.Getenv(Address); addr != "" {
		serverAddr = addr
	}
}

func setupConfiguration() {
	parseFlags()
	parseEnv()
}
