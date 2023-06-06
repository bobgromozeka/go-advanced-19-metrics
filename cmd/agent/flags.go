package main

import "flag"

var serverAddr string
var pollInterval int
var reportInterval int

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "server address to send metrics")
	flag.IntVar(&pollInterval, "p", 2, "Metrics polling interval")
	flag.IntVar(&reportInterval, "r", 10, "Metrics reporting interval to server")

	flag.Parse()
}
