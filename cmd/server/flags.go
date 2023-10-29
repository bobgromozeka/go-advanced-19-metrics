package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bobgromozeka/metrics/internal/server"
)

var startupConfig server.StartupConfig

const JSONConfigPath = "CONFIG"

const (
	Address         = "ADDRESS"
	StoreInterval   = "STORE_INTERVAL"
	FileStoragePath = "FILE_STORAGE_PATH"
	Restore         = "RESTORE"
	DatabaseDsn     = "DATABASE_DSN"
	Key             = "KEY"
	PrivateKeyPath  = "CRYPTO_KEY"
	TrustedSubnet   = "TRUSTED_SUBNET"
)

func parseFlags() {
	flag.StringVar(&startupConfig.ServerAddr, "a", ":8080", "address and port to run server")
	flag.UintVar(&startupConfig.StoreInterval, "i", 300, "Interval of storing metrics to file")
	flag.StringVar(&startupConfig.FileStoragePath, "f", "/tmp/metrics-db.json", "Metrics file storage path")
	flag.BoolVar(&startupConfig.Restore, "r", true, "Restore metrics from file on server start or not")
	flag.StringVar(
		&startupConfig.DatabaseDsn, "d", "",
		"Postgresql data source name (connection string like postgres://practicum:practicum@localhost:5432/practicum)",
	)
	flag.StringVar(&startupConfig.HashKey, "k", "", "Key to validate requests and sign responses")
	flag.StringVar(&startupConfig.PrivateKeyPath, "ck", "./private.pem", "Private key for data encryption")
	flag.StringVar(&startupConfig.TrustedSubnet, "t", "", "Subnet to permit requests from")

	flag.Parse()
}

func parseEnv() {
	if addr := os.Getenv(Address); addr != "" {
		startupConfig.ServerAddr = addr
	}

	if interval := os.Getenv(StoreInterval); interval != "" {
		parsedInterval, err := strconv.Atoi(interval)
		if err != nil {
			log.Fatalln(StoreInterval+" parsing error ", err)
		}
		if parsedInterval < 0 {
			log.Fatalln(StoreInterval + " must be greater or equal 0")
		}
		startupConfig.StoreInterval = uint(parsedInterval)
	}

	if path := os.Getenv(FileStoragePath); path != "" {
		startupConfig.FileStoragePath = path
	}

	if r := os.Getenv(Restore); r == "false" || r == "0" {
		startupConfig.Restore = false
	}

	if dsn := os.Getenv(DatabaseDsn); dsn != "" {
		startupConfig.DatabaseDsn = dsn
	}

	if key := os.Getenv(Key); key != "" {
		startupConfig.HashKey = key
	}

	if privateKeyPath := os.Getenv(PrivateKeyPath); privateKeyPath != "" {
		startupConfig.PrivateKeyPath = privateKeyPath
	}

	if ts := os.Getenv(TrustedSubnet); ts != "" {
		startupConfig.TrustedSubnet = ts
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
