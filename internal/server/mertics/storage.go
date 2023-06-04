package mertics

type metrics struct {
	gauge   float64
	counter int64
}

type MetricsStorage interface {
	UpdateMetricsType(string, string, string) bool
	GetAll() map[string]metrics
}

type MemStorage struct {
	metrics map[string]metrics
}

func New() MemStorage {
	return MemStorage{metrics: map[string]metrics{}}
}

func (s *MemStorage) UpdateMetricsType(metricsType string, name string, value string) bool {
	if _, ok := s.metrics[name]; !ok {
		s.metrics[name] = metrics{}
	}

	switch metricsType {
	case Counter:
		return s.addCounter(name, value)
	case Gauge:
		return s.setGauge(name, value)
	default:
		return false
	}
}

func (s *MemStorage) GetAll() map[string]metrics {
	return s.metrics
}

func (s *MemStorage) addCounter(name string, value string) bool {
	parsedValue, err := parseCounter(value)
	if err != nil {
		return false
	}

	m := s.metrics[name]
	m.counter += parsedValue
	s.metrics[name] = m

	return true
}

func (s *MemStorage) setGauge(name string, value string) bool {
	parsedValue, err := parseGauge(value)
	if err != nil {
		return false
	}

	m := s.metrics[name]
	m.gauge = parsedValue
	s.metrics[name] = m

	return true
}
