package agent

import (
	"log"
	"os"
	"sync"
	"time"
)

var serverAddr string

// Run Starts agent metrics collection and reporting to server.
func Run(c StartupConfig) {
	serverAddr = c.ServerScheme + "://" + c.ServerAddr

	wg := sync.WaitGroup{}
	wg.Add(2)

	rmChan := make(chan runtimeMetrics, 1)
	publicKey, err := os.ReadFile(c.PublicKeyPath)
	if err != nil {
		log.Fatalf("Could not open public key file: %v", err)
	}

	//TODO Add context for graceful shutdown of agent
	go runCollecting(rmChan, c.PollInterval)
	go runReporting(rmChan, c.HashKey, publicKey, c.ReportInterval)

	wg.Wait()
}

func runCollecting(c chan runtimeMetrics, pollInterval int) {
	for {
		rm, err := getRuntimeMetrics()
		if err != nil {
			log.Println("Error during metrics collection: ", err)
		}

		//Clear channel if is not empty to hold the latest value
		if len(c) > 0 {
			<-c
		}

		c <- rm
		time.Sleep(time.Second * time.Duration(pollInterval))
	}
}

func runReporting(c chan runtimeMetrics, hashKey string, publicKey []byte, reportInterval int) {
	for {
		reportToServer(serverAddr, hashKey, publicKey, <-c)
		time.Sleep(time.Second * time.Duration(reportInterval))
	}
}
