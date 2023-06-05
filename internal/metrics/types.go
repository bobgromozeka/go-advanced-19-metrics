package metrics

import "strconv"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var validNames = map[string]struct{}{
	Gauge:   {},
	Counter: {},
}

func ParseCounter(value string) (int64, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func ParseGauge(value string) (float64, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func IsValidValue(metricsType string, value string) bool {
	isValid := false

	switch metricsType {
	case Counter:
		_, err := ParseCounter(value)
		isValid = err == nil
	case Gauge:
		_, err := ParseGauge(value)
		isValid = err == nil
	}

	return isValid
}

func IsValidType(metricsType string) bool {
	_, ok := validNames[metricsType]
	return ok
}
