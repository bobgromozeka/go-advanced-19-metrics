package storage

import (
	"encoding/json"
	"os"

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
}

type Event int

const (
	Update Event = iota
)

type MemStorage struct {
	Metrics   Metrics
	Listeners map[Event][]func()
}

type Metrics struct {
	gauge   GaugeMetrics
	counter CounterMetrics
}

func New() MemStorage {
	return MemStorage{
		Metrics: Metrics{
			gauge:   GaugeMetrics{},
			counter: CounterMetrics{},
		},
		Listeners: map[Event][]func(){},
	}
}

func (s MemStorage) GetMetricsByType(mtype string, name string) (v any, ok bool) {
	switch mtype {
	case metrics.GaugeType:
		return s.GetGaugeMetrics(name)
	case metrics.CounterType:
		return s.GetCounterMetrics(name)
	default:
		return nil, false
	}
}

func (s MemStorage) GetAllGaugeMetrics() GaugeMetrics {
	return s.Metrics.gauge
}

func (s MemStorage) GetAllCounterMetrics() CounterMetrics {
	return s.Metrics.counter
}

func (s MemStorage) GetGaugeMetrics(name string) (v float64, ok bool) {
	v, ok = s.GetAllGaugeMetrics()[name]
	return v, ok
}

func (s MemStorage) GetCounterMetrics(name string) (v int64, ok bool) {
	v, ok = s.GetAllCounterMetrics()[name]
	return v, ok
}

func (s MemStorage) AddCounter(name string, value int64) int64 {
	if _, ok := s.Metrics.counter[name]; !ok {
		s.Metrics.counter[name] = 0
	}

	s.Metrics.counter[name] += value

	s.fireEvent(Update)

	return s.Metrics.counter[name]
}

func (s MemStorage) SetGauge(name string, value float64) float64 {
	if _, ok := s.Metrics.gauge[name]; !ok {
		s.Metrics.gauge[name] = 0
	}

	s.Metrics.gauge[name] = value

	s.fireEvent(Update)

	return value
}

func (s MemStorage) UpdateMetricsByType(metricsType string, name string, value string) (any, error) {
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

func (s MemStorage) Listen(event Event, callbacks ...func()) {
	s.Listeners[event] = append(s.Listeners[event], callbacks...)
}

func (s MemStorage) fireEvent(event Event) {
	for _, callback := range s.Listeners[event] {
		callback()
	}
}

func (s MemStorage) PersistToPath(path string) error {
	dataToMarshal := struct {
		Counter CounterMetrics
		Gauge   GaugeMetrics
	}{
		Counter: s.Metrics.counter,
		Gauge:   s.Metrics.gauge,
	}
	jsonData, jsonErr := json.Marshal(dataToMarshal)
	if jsonErr != nil {
		return jsonErr
	}

	if writeErr := os.WriteFile(path, jsonData, 0666); writeErr != nil {
		return writeErr
	}

	return nil
}

func (s MemStorage) RestoreFrom(filepath string) error {
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	unmarshalledData := struct {
		Counter CounterMetrics
		Gauge   GaugeMetrics
	}{}

	unmarshalErr := json.Unmarshal(jsonData, &unmarshalledData)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	for k, v := range unmarshalledData.Gauge {
		s.Metrics.gauge[k] = v
	}

	for k, v := range unmarshalledData.Counter {
		s.Metrics.counter[k] = v
	}

	return nil
}
