package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillRuntimeMetricsUpdatesPollCount(t *testing.T) {
	rm, _ := getRuntimeMetrics()

	assert.EqualValues(t, 1, rm.PollCount)
}

func TestFillRuntimeMetricsUpdatesRandomValue(t *testing.T) {
	rm, _ := getRuntimeMetrics()

	assert.NotEqualValues(t, 0, rm.RandomValue)
}
