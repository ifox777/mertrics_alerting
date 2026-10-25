package repository

import "sync"

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

type Storage interface {
	SetGauge(id string, value float64)
	AddCounter(id string, delta int64)
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
