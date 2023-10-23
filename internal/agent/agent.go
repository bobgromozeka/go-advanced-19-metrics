package agent

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

var serverAddr string

// Run Starts agent metrics collection and reporting to server.
func Run(ctx context.Context, c StartupConfig) {
	serverAddr = c.ServerScheme + "://" + c.ServerAddr

	wg := sync.WaitGroup{}

	rmChan := make(chan runtimeMetrics, 1)
	publicKey, err := os.ReadFile(c.PublicKeyPath)
	if err != nil {
		log.Fatalf("Could not open public key file: %v", err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		runCollecting(ctx, rmChan, c.PollInterval)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		runReporting(ctx, rmChan, c.HashKey, publicKey, c.ReportInterval)
	}()

	wg.Wait()
}

func runCollecting(ctx context.Context, c chan runtimeMetrics, pollInterval int) {
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

		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Second * time.Duration(pollInterval))
	}
}

func runReporting(ctx context.Context, c chan runtimeMetrics, hashKey string, publicKey []byte, reportInterval int) {
	for {
		reportToServer(serverAddr, hashKey, publicKey, <-c)

		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Second * time.Duration(reportInterval))
	}
}
