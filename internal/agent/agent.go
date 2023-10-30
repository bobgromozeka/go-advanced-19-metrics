package agent

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

// Run Starts agent metrics collection and reporting to server.
func Run(ctx context.Context, c StartupConfig) {

	wg := sync.WaitGroup{}

	rmChan := make(chan runtimeMetrics, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()
		runCollecting(ctx, rmChan, c.PollInterval)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		if c.ReportGRPC {
			serverAddr := c.ServerAddr
			runReportingGRPC(ctx, serverAddr, rmChan, c.PublicKeyPath, c.ReportInterval)
		} else {
			serverAddr := c.ServerScheme + "://" + c.ServerAddr
			runReportingHTTP(ctx, serverAddr, rmChan, c.HashKey, c.PublicKeyPath, c.ReportInterval)
		}
	}()

	log.Println("Started polling and sending data to server.....")

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

func runReportingGRPC(ctx context.Context, serverAddr string, c chan runtimeMetrics, certPath string, reportInterval int) {
	for {
		reportToGRPCServer(serverAddr, certPath, <-c)

		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Second * time.Duration(reportInterval))
	}
}

func runReportingHTTP(ctx context.Context, serverAddr string, c chan runtimeMetrics, hashKey string, publicKeyPath string, reportInterval int) {
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatalf("Could not open public key file: %v", err)
	}

	for {
		reportToHTTPServer(serverAddr, hashKey, publicKey, <-c)

		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Second * time.Duration(reportInterval))
	}
}
