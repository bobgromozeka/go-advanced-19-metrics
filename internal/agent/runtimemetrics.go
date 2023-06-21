package agent

import (
	"math/rand"
	"runtime"

	"github.com/bobgromozeka/metrics/internal/metrics"
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

func fillRuntimeMetrics(rm *runtimeMetrics) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	rnd := rand.Float64()

	rm.Alloc = ms.Alloc
	rm.BuckHashSys = ms.BuckHashSys
	rm.Frees = ms.Frees
	rm.GCCPUFraction = ms.GCCPUFraction
	rm.GCSys = ms.GCSys
	rm.HeapAlloc = ms.HeapAlloc
	rm.HeapIdle = ms.HeapIdle
	rm.HeapInuse = ms.HeapInuse
	rm.HeapObjects = ms.HeapObjects
	rm.HeapReleased = ms.HeapReleased
	rm.HeapSys = ms.HeapSys
	rm.LastGC = ms.LastGC
	rm.Lookups = ms.Lookups
	rm.MCacheInuse = ms.MCacheInuse
	rm.MCacheSys = ms.MCacheSys
	rm.MSpanInuse = ms.MSpanInuse
	rm.MSpanSys = ms.MSpanSys
	rm.Mallocs = ms.Mallocs
	rm.NextGC = ms.NextGC
	rm.NumForcedGC = ms.NumForcedGC
	rm.NumGC = ms.NumGC
	rm.OtherSys = ms.OtherSys
	rm.PauseTotalNs = ms.PauseTotalNs
	rm.StackInuse = ms.StackInuse
	rm.StackSys = ms.StackSys
	rm.Sys = ms.Sys
	rm.TotalAlloc = ms.TotalAlloc
	rm.PollCount++
	rm.RandomValue = rnd * 1000
}
