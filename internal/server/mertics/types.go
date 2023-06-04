package mertics

import "strconv"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var validNames = map[string]struct{}{
	Gauge:   {},
	Counter: {},
}

func parseCounter(value string) (int64, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func parseGauge(value string) (float64, error) {
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
		_, err := parseCounter(value)
		isValid = err == nil
	case Gauge:
		_, err := parseGauge(value)
		isValid = err == nil
	}

	return isValid
}

func IsValidType(metricsType string) bool {
	_, ok := validNames[metricsType]
	return ok
}
