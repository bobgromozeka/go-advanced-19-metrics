package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var serverAddr string
var storeInterval uint
var fileStoragePath string
var restore bool

const (
	Address         = "ADDRESS"
	StoreInterval   = "STORE_INTERVAL"
	FileStoragePath = "FILE_STORAGE_PATH"
	Restore         = "RESTORE"
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", ":8080", "address and port to run server")
	flag.UintVar(&storeInterval, "i", 300, "Interval of storing metrics to file")
	flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "Metrics file storage path")
	flag.BoolVar(&restore, "r", true, "Restore metrics from file on server start or not")

	flag.Parse()
}

func parseEnv() {
	if addr := os.Getenv(Address); addr != "" {
		serverAddr = addr
	}
	if interval := os.Getenv(StoreInterval); interval != "" {
		parsedInterval, err := strconv.Atoi(interval)
		if err != nil {
			log.Fatalln(StoreInterval+" parsing error ", err)
		}
		if parsedInterval < 0 {
			log.Fatalln(StoreInterval + " must be greater or equal 0")
		}
		storeInterval = uint(parsedInterval)
	}
	if path := os.Getenv(FileStoragePath); path != "" {
		fileStoragePath = path
	}
	if r := os.Getenv(Restore); r == "false" || r == "0" {
		restore = false
	}
}

func setupConfiguration() {
	parseFlags()
	parseEnv()
}
