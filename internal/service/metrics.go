package service

import (
	"fmt"
	"strconv"
)

type Metrics struct {
	storage Storage
}

func NewMetrics(storage Storage) *Metrics {
	return &Metrics{
		storage: storage,
	}
}

func (m *Metrics) Update(kind, name, rawValue string) error {
	switch kind {
	case "gauge":
		gauge, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			return err
		}
		m.storage.SetGauge(name, gauge)
	case "counter":
		counter, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			return err
		}
		m.storage.AddCounter(name, counter)
	default:

		return fmt.Errorf("unknown metric type %q", kind)
	}
	return nil
}

func (m *Metrics) Value(kind, name string) (string, error) {
	switch kind {
	case "gauge":
		value, ok := m.storage.GetGauge(name)
		if !ok {
			return "", fmt.Errorf("gauge metric %q not found", name)
		}
		return fmt.Sprint(value), nil
	case "counter":
		value, ok := m.storage.GetCounter(name)
		if !ok {
			return "", fmt.Errorf("counter metric %q not found", name)
		}
		return fmt.Sprint(value), nil
	default:
		return "", fmt.Errorf("unknown metric type %q", kind)
	}
}

func (m *Metrics) GetAll() (map[string]float64, map[string]int64) {
	return m.storage.GetAll()
}
