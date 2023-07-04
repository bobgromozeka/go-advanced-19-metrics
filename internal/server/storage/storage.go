package storage

import (
	"github.com/bobgromozeka/metrics/internal/metrics"
)

type GaugeMetrics map[string]float64
type CounterMetrics map[string]int64

type Storage interface {
	AddCounter(string, int64) int64
	SetGauge(string, float64) float64
	UpdateMetricsByType(string, string, string) (any, error)
	GetAllGaugeMetrics() GaugeMetrics
	GetAllCounterMetrics() CounterMetrics
	GetGaugeMetrics(string) (float64, bool)
	GetCounterMetrics(string) (int64, bool)
	GetMetricsByType(string, string) (any, bool)
	SetMetrics(Metrics)
}

type MemStorage struct {
	Metrics Metrics
}

type Metrics struct {
	Gauge   GaugeMetrics
	Counter CounterMetrics
}

func New() Storage {
	return &MemStorage{
		Metrics: Metrics{
			Gauge:   GaugeMetrics{},
			Counter: CounterMetrics{},
		},
	}
}

func (s *MemStorage) SetMetrics(m Metrics) {
	s.Metrics = m
}

func (s *MemStorage) GetMetricsByType(mtype string, name string) (v any, ok bool) {
	switch mtype {
	case metrics.GaugeType:
		return s.GetGaugeMetrics(name)
	case metrics.CounterType:
		return s.GetCounterMetrics(name)
	default:
		return nil, false
	}
}

func (s *MemStorage) GetAllGaugeMetrics() GaugeMetrics {
	return s.Metrics.Gauge
}

func (s *MemStorage) GetAllCounterMetrics() CounterMetrics {
	return s.Metrics.Counter
}

func (s *MemStorage) GetGaugeMetrics(name string) (v float64, ok bool) {
	v, ok = s.GetAllGaugeMetrics()[name]
	return v, ok
}

func (s *MemStorage) GetCounterMetrics(name string) (v int64, ok bool) {
	v, ok = s.GetAllCounterMetrics()[name]
	return v, ok
}

func (s *MemStorage) AddCounter(name string, value int64) int64 {
	if _, ok := s.Metrics.Counter[name]; !ok {
		s.Metrics.Counter[name] = 0
	}

	s.Metrics.Counter[name] += value

	return s.Metrics.Counter[name]
}

func (s *MemStorage) SetGauge(name string, value float64) float64 {
	if _, ok := s.Metrics.Gauge[name]; !ok {
		s.Metrics.Gauge[name] = 0
	}

	s.Metrics.Gauge[name] = value

	return value
}

func (s *MemStorage) UpdateMetricsByType(metricsType string, name string, value string) (any, error) {
	switch metricsType {
	case metrics.CounterType:
		parsedValue, err := metrics.ParseCounter(value)
		if err != nil {
			return false, err
		}
		return s.AddCounter(name, parsedValue), nil
	case metrics.GaugeType:
		parsedValue, err := metrics.ParseGauge(value)
		if err != nil {
			return false, err
		}
		return s.SetGauge(name, parsedValue), nil
	default:
		return false, nil
	}
}
