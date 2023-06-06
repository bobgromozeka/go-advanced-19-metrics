package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidType(t *testing.T) {
	tests := []struct {
		name string
		Type string
		want bool
	}{
		{
			name: "GaugeType is valid",
			Type: "gauge",
			want: true,
		},
		{
			name: "CounterType is valid",
			Type: "counter",
			want: true,
		},
		{
			name: "Random is valid",
			Type: "random",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidType(tt.Type))
		})
	}
}

func TestIsValidValue(t *testing.T) {
	type args struct {
		metricsType string
		value       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Number is valid for gauge",
			args: args{
				metricsType: "gauge",
				value:       "123",
			},
			want: true,
		},
		{
			name: "Not Number is not valid for gauge",
			args: args{
				metricsType: "gauge",
				value:       "1s23r",
			},
			want: false,
		},
		{
			name: "Number is valid for counter",
			args: args{
				metricsType: "gauge",
				value:       "123",
			},
			want: true,
		},
		{
			name: "Not Number is not valid for counter",
			args: args{
				metricsType: "gauge",
				value:       "1s23r",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidValue(tt.args.metricsType, tt.args.value))
		})
	}
}
