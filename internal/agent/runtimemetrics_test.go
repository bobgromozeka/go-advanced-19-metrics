package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillRuntimeMetricsUpdatesPollCount(t *testing.T) {
	rm := runtimeMetrics{}

	fillRuntimeMetrics(&rm)

	assert.EqualValues(t, 1, rm.PollCount)
}

func TestFillRuntimeMetricsUpdatesRandomValue(t *testing.T) {
	rm := runtimeMetrics{}

	fillRuntimeMetrics(&rm)

	assert.NotEqualValues(t, 0, rm.RandomValue)
}
