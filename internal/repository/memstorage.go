package repository

import (
	"maps"
	"sync"
)

type MemStorage struct {
	mu       sync.Mutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MemStorage) SetGauge(id string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[id] = value
}

func (m *MemStorage) AddCounter(id string, delta int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[id] += delta
}

func (m *MemStorage) GetGauge(id string) (float64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.gauges[id]
	return value, ok
}

func (m *MemStorage) GetCounter(id string) (int64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, ok := m.counters[id]
	return value, ok
}

func (m *MemStorage) GetAll() (map[string]float64, map[string]int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	gauses := maps.Clone(m.gauges)
	counters := maps.Clone(m.counters)

	return gauses, counters
}
