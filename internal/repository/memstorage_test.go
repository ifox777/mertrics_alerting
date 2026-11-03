package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetGauge(t *testing.T) {
	memStorage := NewMemStorage()

	memStorage.SetGauge("temperature", 21.5)
	memStorage.SetGauge("temperature", 22.5)
	assert.Equal(t, 22.5, memStorage.gauges["temperature"])
}

func TestMemStorage_SetCounter(t *testing.T) {
	memStorage := NewMemStorage()

	memStorage.AddCounter("PollCount", 2)
	memStorage.AddCounter("PollCount", 3)

	assert.Equal(t, int64(5), memStorage.counters["PollCount"])
}
