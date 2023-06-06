package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_makeEndpointsFromStructure(t *testing.T) {
	type args struct {
		rm any
	}
	tests := []struct {
		name     string
		args     args
		want     []string
		positive bool
	}{
		{
			name: `Makes "counter" update from PollCount`,
			args: args{
				rm: runtimeMetrics{
					PollCount: 123,
				},
			},
			want: []string{
				"/update/counter/PollCount/123",
			},
			positive: true,
		},
		{
			name: "Can work with uint64 and uint32",
			args: args{
				rm: runtimeMetrics{
					Alloc: 1,
					NumGC: 2,
				},
			},
			want: []string{
				"/update/gauge/Alloc/1",
				"/update/gauge/NumGC/2",
			},
			positive: true,
		},
		{
			name: "Skips unknown types",
			args: args{
				rm: struct {
					RandomField string
				}{
					RandomField: "asd",
				},
			},
			want: []string{
				"/update/gauge/RandomField/asd",
			},
			positive: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.positive {
				assert.Subset(t, makeEndpointsFromStructure(test.args.rm), test.want)
			} else {
				assert.NotSubset(t, makeEndpointsFromStructure(test.args.rm), test.want)
			}
		})
	}
}
