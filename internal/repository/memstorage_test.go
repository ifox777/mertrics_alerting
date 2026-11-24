package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestMemStorage_GetGauge(t *testing.T) {
	memStorage := NewMemStorage()
	memStorage.SetGauge("temperature", 21.5)
	memStorage.SetGauge("zero", 0)

	tests := []struct {
		name  string
		id    string
		value float64
		has   bool
	}{
		{"saved 21.5", "temperature", float64(21.5), true},
		{"saved 0", "zero", float64(0), true},
		{"No name", "missing", float64(0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, has := memStorage.GetGauge(tt.id)

			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.has, has)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	memStorage := NewMemStorage()
	memStorage.AddCounter("PollCount", 2)
	memStorage.AddCounter("PollCount", 3)
	memStorage.AddCounter("zero", 0)

	tests := []struct {
		name  string
		id    string
		value int64
		has   bool
	}{
		{"saved sum", "PollCount", 5, true},
		{"saved zero", "zero", 0, true},
		{"missing", "missing", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, has := memStorage.GetCounter(tt.id)

			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.has, has)
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	storage := NewMemStorage()
	storage.SetGauge("gauge", 21.5)
	storage.AddCounter("counter", 5)

	gauges, counters := storage.GetAll()

	require.Equal(t, map[string]float64{"gauge": 21.5}, gauges)
	require.Equal(t, map[string]int64{"counter": 5}, counters)

	gauges["gauge"] = 999
	counters["counter"] = 1000
	gaugeValue, gaugeExists := storage.GetGauge("gauge")
	counterValue, counterExists := storage.GetCounter("counter")

	assert.True(t, gaugeExists)
	assert.True(t, counterExists)
	assert.Equal(t, 21.5, gaugeValue)
	assert.Equal(t, int64(5), counterValue)
}

func TestSaveToFileCounter(t *testing.T) {
	storage := NewMemStorage()
	storage.AddCounter("PollCount", 2)
	storage.AddCounter("PollCount", 3)

	path := filepath.Join(t.TempDir(), "metrics.json")
	require.NoError(t, storage.SaveToFile(path))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.JSONEq(t, `[{"id":"PollCount","type":"counter","delta":5}]`, string(data))
}

func TestSaveAndLoadMetrics(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metrics.json")
	source := NewMemStorage()
	source.SetGauge("Alloc", 0)
	source.AddCounter("PollCount", 5)
	require.NoError(t, source.SaveToFile(path))

	restored := NewMemStorage()
	require.NoError(t, restored.LoadFromFile(path))
	require.NoError(t, restored.LoadFromFile(path))

	gauges, counters := restored.GetAll()
	require.Equal(t, map[string]float64{"Alloc": 0}, gauges)
	require.Equal(t, map[string]int64{"PollCount": 5}, counters)
}
