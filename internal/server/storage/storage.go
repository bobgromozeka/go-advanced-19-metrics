package storage

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("metrics not found")
var ErrWrongMetrics = errors.New("wrong metrics type")

type GaugeMetrics map[string]float64
type CounterMetrics map[string]int64

// PersistenceSettings
// Path File path. Note that PersistentStorage saves data in JSON format.
// Interval Seconds between saving data into file. If 0 - saves data synchronously on every update.
// Restore Load saved data on application start or not
type PersistenceSettings struct {
	Path     string
	Interval uint
	Restore  bool
}

type Storage interface {
	// AddCounter Adds new counter or increases existing one by name.
	AddCounter(context.Context, string, int64) (int64, error)
	// AddCounters Adds new counters or increases existing by names.
	AddCounters(context.Context, CounterMetrics) error
	//SetGauge Sets gauge metric value by name
	SetGauge(context.Context, string, float64) (float64, error)
	//SetGauges Sets gauge metrics values by name
	SetGauges(context.Context, GaugeMetrics) error
	//UpdateMetricsByType Sets gauge or adds counter depending on metrics type
	UpdateMetricsByType(context.Context, string, string, string) (any, error)
	// GetAllGaugeMetrics Returns all gauge metrics from storage
	GetAllGaugeMetrics(context.Context) (GaugeMetrics, error)
	// GetAllCounterMetrics Returns all counter metrics from storage
	GetAllCounterMetrics(context.Context) (CounterMetrics, error)
	// GetGaugeMetrics Returns gauge metrics by name
	GetGaugeMetrics(context.Context, string) (float64, error)
	// GetCounterMetrics Returns counter metrics by name
	GetCounterMetrics(context.Context, string) (int64, error)
	// GetMetricsByType Returns metrics depending on type
	GetMetricsByType(context.Context, string, string) (any, error)
}
