package storage

import (
	"github.com/bobgromozeka/metrics/internal/metrics"
)

type GaugeMetrics map[string]float64
type CounterMetrics map[string]int64

type Storage interface {
	UpdateMetricsType(string, string, string) (bool, error)
	GetAllGaugeMetrics() GaugeMetrics
	GetAllCounterMetrics() CounterMetrics
	GetMetrics(metricsType string, name string) (any, bool)
}

type MemStorage struct {
	gaugeMetrics   GaugeMetrics
	counterMetrics CounterMetrics
}

func New() MemStorage {
	return MemStorage{
		gaugeMetrics:   GaugeMetrics{},
		counterMetrics: CounterMetrics{},
	}
}

func (s MemStorage) UpdateMetricsType(metricsType string, name string, value string) (bool, error) {
	switch metricsType {
	case metrics.CounterType:
		return s.addCounter(name, value)
	case metrics.GaugeType:
		return s.setGauge(name, value)
	default:
		return false, nil
	}
}

func (s MemStorage) GetAllGaugeMetrics() GaugeMetrics {
	return s.gaugeMetrics
}

func (s MemStorage) GetAllCounterMetrics() CounterMetrics {
	return s.counterMetrics
}

func (s MemStorage) GetMetrics(metricsType string, name string) (v any, ok bool) {
	switch metricsType {
	case metrics.GaugeType:
		v, ok = s.GetAllGaugeMetrics()[name]
	case metrics.CounterType:
		v, ok = s.GetAllCounterMetrics()[name]
	}

	return
}

func (s MemStorage) addCounter(name string, value string) (bool, error) {
	parsedValue, err := metrics.ParseCounter(value)
	if err != nil {
		return false, err
	}

	if _, ok := s.counterMetrics[name]; !ok {
		s.counterMetrics[name] = 0
	}

	s.counterMetrics[name] += parsedValue

	return true, nil
}

func (s MemStorage) setGauge(name string, value string) (bool, error) {
	parsedValue, err := metrics.ParseGauge(value)
	if err != nil {
		return false, err
	}

	if _, ok := s.gaugeMetrics[name]; !ok {
		s.gaugeMetrics[name] = 0
	}

	s.gaugeMetrics[name] = parsedValue

	return true, nil
}
