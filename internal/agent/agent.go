package agent

import (
	"github.com/bobgromozeka/metrics/internal/metrics"
	"sync"
	"time"
)

var runtimeMetricsTypes = map[string]string{
	//GaugeType if not specified here
	"PollCount": metrics.CounterType,
}

type runtimeMetrics struct {
	Alloc         uint64
	TotalAlloc    uint64
	Sys           uint64
	Lookups       uint64
	Mallocs       uint64
	Frees         uint64
	HeapAlloc     uint64
	HeapSys       uint64
	HeapIdle      uint64
	HeapInuse     uint64
	HeapReleased  uint64
	HeapObjects   uint64
	StackInuse    uint64
	StackSys      uint64
	MSpanInuse    uint64
	MSpanSys      uint64
	MCacheInuse   uint64
	MCacheSys     uint64
	BuckHashSys   uint64
	GCSys         uint64
	OtherSys      uint64
	NextGC        uint64
	LastGC        uint64
	PauseTotalNs  uint64
	NumGC         uint32
	NumForcedGC   uint32
	GCCPUFraction float64
	PollCount     uint64
	RandomValue   float64
}

func Run() {
	rm := runtimeMetrics{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for {
			fillRuntimeMetrics(&rm)
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 10)
			reportToServer(rm) // Wait while some metrics are collected
		}
	}()

	wg.Wait()
}
