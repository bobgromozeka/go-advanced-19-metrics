package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStorage_GetAll(t *testing.T) {
	s := MemStorage{
		metrics: map[string]any{
			"name": float64(1),
		},
	}

	content := s.GetAll()

	require.Contains(t, content, "name")
	assert.EqualValues(t, 1, content["name"])
}

func TestMemStorage_addCounter(t *testing.T) {
	type fields struct {
		metrics map[string]any
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
				metrics: map[string]any{
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
				metrics: map[string]any{},
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
				metrics: map[string]any{
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
				metrics: tt.fields.metrics,
			}
			assert.Equal(t, tt.want, s.addCounter(tt.args.name, tt.args.value))
			require.Contains(t, s.metrics, tt.args.name)
			assert.EqualValues(t, tt.wantValue, s.metrics[tt.args.name])
		})
	}
}

func TestMemStorage_setGauge(t *testing.T) {
	type fields struct {
		metrics map[string]any
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
				metrics: map[string]any{
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
				metrics: map[string]any{},
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
				metrics: map[string]any{
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
				metrics: tt.fields.metrics,
			}
			if got := s.setGauge(tt.args.name, tt.args.value); got != tt.want {
				assert.Equal(t, tt.want, s.addCounter(tt.args.name, tt.args.value))
				require.Contains(t, s.metrics, tt.args.name)
				assert.EqualValues(t, tt.wantValue, s.metrics[tt.args.name])
			}
		})
	}
}
