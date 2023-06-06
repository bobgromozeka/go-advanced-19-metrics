package main

import "flag"

var serverAddr string

func parseFlags() {
	flag.StringVar(&serverAddr, "a", ":8080", "address and port to run server")

	flag.Parse()
}
