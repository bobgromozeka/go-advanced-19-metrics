package storage

import (
	"github.com/bobgromozeka/metrics/internal/metrics"
)

type Storage interface {
	UpdateMetricsType(string, string, string) bool
	GetAll() map[string]any
}

type MemStorage struct {
	metrics map[string]any
}

func New() MemStorage {
	return MemStorage{metrics: map[string]any{}}
}

func (s MemStorage) UpdateMetricsType(metricsType string, name string, value string) bool {
	switch metricsType {
	case metrics.Counter:
		return s.addCounter(name, value)
	case metrics.Gauge:
		return s.setGauge(name, value)
	default:
		return false
	}
}

func (s MemStorage) GetAll() map[string]any {
	return s.metrics
}

func (s MemStorage) addCounter(name string, value string) bool {
	parsedValue, err := metrics.ParseCounter(value)
	if err != nil {
		return false
	}

	if _, ok := s.metrics[name]; !ok {
		s.metrics[name] = int64(0)
	}

	if v, ok := s.metrics[name].(int64); ok {
		s.metrics[name] = v + parsedValue
	} else {
		return false
	}

	return true
}

func (s MemStorage) setGauge(name string, value string) bool {
	parsedValue, err := metrics.ParseGauge(value)
	if err != nil {
		return false
	}

	if _, ok := s.metrics[name]; !ok {
		s.metrics[name] = float64(0)
	}

	if _, ok := s.metrics[name].(float64); ok {
		s.metrics[name] = parsedValue
	} else {
		return false
	}

	return true
}
