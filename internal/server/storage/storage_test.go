package storage

import (
	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStorage_GetAllGaugeMetrics(t *testing.T) {
	s := MemStorage{
		gaugeMetrics: map[string]metrics.Gauge{
			"name": metrics.Gauge(1),
		},
	}

	content := s.GetAllGaugeMetrics()

	require.Contains(t, content, "name")
	assert.EqualValues(t, 1, content["name"])
}

func TestMemStorage_GetAllCounterMetrics(t *testing.T) {
	s := MemStorage{
		counterMetrics: map[string]metrics.Counter{
			"name": metrics.Counter(1),
		},
	}

	content := s.GetAllCounterMetrics()

	require.Contains(t, content, "name")
	assert.EqualValues(t, 1, content["name"])
}

func TestMemStorage_addCounter(t *testing.T) {
	type fields struct {
		metrics map[string]metrics.Counter
	}
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		wantValue uint64
	}{
		{
			name: "successfully adds value to counter",
			fields: fields{
				metrics: map[string]metrics.Counter{
					"c": int64(1),
				},
			},
			args: args{
				name:  "c",
				value: "1",
			},
			want:      true,
			wantValue: 2,
		},
		{
			name: "successfully creates new metrics",
			fields: fields{
				metrics: map[string]metrics.Counter{},
			},
			args: args{
				name:  "c",
				value: "1",
			},
			want:      true,
			wantValue: 1,
		},
		{
			name: "ignores wrong values",
			fields: fields{
				metrics: map[string]metrics.Counter{
					"c": int64(1),
				},
			},
			args: args{
				name:  "c",
				value: "1wrongvalue",
			},
			want:      false,
			wantValue: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MemStorage{
				counterMetrics: tt.fields.metrics,
			}
			added, _ := s.addCounter(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want, added)
			require.Contains(t, s.counterMetrics, tt.args.name)
			assert.EqualValues(t, tt.wantValue, s.counterMetrics[tt.args.name])
		})
	}
}

func TestMemStorage_setGauge(t *testing.T) {
	type fields struct {
		metrics map[string]metrics.Gauge
	}
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		wantValue uint64
	}{
		{
			name: "successfully sets value to gauge",
			fields: fields{
				metrics: map[string]metrics.Gauge{
					"c": float64(1),
				},
			},
			args: args{
				name:  "c",
				value: "123",
			},
			want:      true,
			wantValue: 123,
		},
		{
			name: "successfully creates new metrics",
			fields: fields{
				metrics: map[string]metrics.Gauge{},
			},
			args: args{
				name:  "c",
				value: "122",
			},
			want:      true,
			wantValue: 122,
		},
		{
			name: "ignores wrong values",
			fields: fields{
				metrics: map[string]metrics.Gauge{
					"c": float64(1),
				},
			},
			args: args{
				name:  "c",
				value: "1wrongvalue",
			},
			want:      false,
			wantValue: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MemStorage{
				gaugeMetrics: tt.fields.metrics,
			}
			if got, _ := s.setGauge(tt.args.name, tt.args.value); got != tt.want {
				assert.Equal(t, tt.want, got)
				require.Contains(t, s.gaugeMetrics, tt.args.name)
				assert.EqualValues(t, tt.wantValue, s.gaugeMetrics[tt.args.name])
			}
		})
	}
}
