package agent

import (
	"testing"

	"github.com/bobgromozeka/metrics/internal/metrics"

	"github.com/stretchr/testify/assert"
)

func Test_makeBodiesFromStructure(t *testing.T) {
	counterValue1 := int64(123)
	counterValue2 := float64(1)
	counterValue3 := float64(2)
	type args struct {
		rm any
	}
	tests := []struct {
		name string
		args args
		want []metrics.RequestPayload
	}{
		{
			name: `Makes "counter" update from PollCount`,
			args: args{
				rm: runtimeMetrics{
					PollCount: 123,
				},
			},
			want: []metrics.RequestPayload{
				{
					ID:    "PollCount",
					MType: "counter",
					Delta: &counterValue1,
				},
			},
		},
		{
			name: "Can work with uint64 and uint32",
			args: args{
				rm: runtimeMetrics{
					Alloc: 1,
					NumGC: 2,
				},
			},
			want: []metrics.RequestPayload{
				{
					ID:    "Alloc",
					MType: "gauge",
					Value: &counterValue2,
				},
				{
					ID:    "NumGC",
					MType: "gauge",
					Value: &counterValue3,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Subset(t, makeBodiesFromStructure(test.args.rm), test.want)
		})
	}
}

func Test_makeBodiesSkipsUnknownTypes(t *testing.T) {
	randomRM := struct {
		RandomField string
	}{
		RandomField: "asd",
	}

	assert.Len(t, makeBodiesFromStructure(randomRM), 0)
}
