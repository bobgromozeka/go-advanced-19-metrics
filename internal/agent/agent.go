package agent

import (
	"sync"
	"time"
)

var serverAddr string

func Run(c StartupConfig) {
	serverAddr = c.ServerAddr
	rm := runtimeMetrics{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for {
			fillRuntimeMetrics(&rm)
			time.Sleep(time.Second * time.Duration(c.PollInterval))
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(c.ReportInterval)) // Wait while some metrics are collected
			reportToServer(rm)
		}
	}()

	wg.Wait()
}
