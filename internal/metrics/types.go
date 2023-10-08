package metrics

import (
	"strconv"
)

type Gauge = float64
type Counter = int64

// RequestPayload payload for updating metrics in JSON format.
type RequestPayload struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

// Metrics types
const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

var validNames = map[string]struct{}{
	GaugeType:   {},
	CounterType: {},
}

// ParseCounter parses string into Counter type or error if not possible.
func ParseCounter(value string) (Counter, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

// ParseGauge parses string into Gauge type or error if not possible.
func ParseGauge(value string) (Gauge, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

// IsValidValue Checks if value is valid according to specified metrics type.
func IsValidValue(metricsType string, value string) bool {
	isValid := false

	switch metricsType {
	case CounterType:
		_, err := ParseCounter(value)
		isValid = err == nil
	case GaugeType:
		_, err := ParseGauge(value)
		isValid = err == nil
	}

	return isValid
}

// IsValidType Checks if metrics type is valid
func IsValidType(metricsType string) bool {
	_, ok := validNames[metricsType]
	return ok
}
