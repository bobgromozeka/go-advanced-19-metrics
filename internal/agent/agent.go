package agent

import (
	"log"
	"sync"
	"time"
)

var serverAddr string

func Run(c StartupConfig) {
	serverAddr = c.ServerScheme + "://" + c.ServerAddr

	wg := sync.WaitGroup{}
	wg.Add(2)

	rmChan := make(chan runtimeMetrics, 1)

	//TODO Add context for graceful shutdown of agent
	go runCollecting(rmChan, c.PollInterval)
	go runReporting(rmChan, c.HashKey, c.ReportInterval)

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

func runReporting(c chan runtimeMetrics, hashKey string, reportInterval int) {
	for {
		reportToServer(serverAddr, hashKey, <-c)
		time.Sleep(time.Second * time.Duration(reportInterval))
	}
}
