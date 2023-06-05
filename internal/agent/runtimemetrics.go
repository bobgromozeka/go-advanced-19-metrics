package agent

import (
	"math/rand"
	"runtime"
)

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
