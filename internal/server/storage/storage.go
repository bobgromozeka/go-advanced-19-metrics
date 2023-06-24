package storage

type GaugeMetrics map[string]float64
type CounterMetrics map[string]int64

type Storage interface {
	AddCounter(string, int64) int64
	SetGauge(string, float64) float64
	GetAllGaugeMetrics() GaugeMetrics
	GetAllCounterMetrics() CounterMetrics
	GetGaugeMetrics(string) (float64, bool)
	GetCounterMetrics(string) (int64, bool)
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

func (s MemStorage) GetAllGaugeMetrics() GaugeMetrics {
	return s.gaugeMetrics
}

func (s MemStorage) GetAllCounterMetrics() CounterMetrics {
	return s.counterMetrics
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
	if _, ok := s.counterMetrics[name]; !ok {
		s.counterMetrics[name] = 0
	}

	s.counterMetrics[name] += value

	return s.counterMetrics[name]
}

func (s MemStorage) SetGauge(name string, value float64) float64 {
	if _, ok := s.gaugeMetrics[name]; !ok {
		s.gaugeMetrics[name] = 0
	}

	s.gaugeMetrics[name] = value

	return value
}
