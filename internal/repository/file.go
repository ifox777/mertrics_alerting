package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	models "github.com/ifox777/metrics-alerting/internal/model"
)

func (m *MemStorage) SaveToFile(path string) error {
	m.saveMu.Lock()
	defer m.saveMu.Unlock()

	gauges, counters := m.GetAll()
	metrics := make([]models.Metrics, 0, len(gauges)+len(counters))

	for id, value := range gauges {
		metrics = append(metrics, models.Metrics{
			ID:    id,
			MType: models.Gauge,
			Value: &value,
		})
	}
	for id, total := range counters {
		metrics = append(metrics, models.Metrics{
			ID:    id,
			MType: models.Counter,
			Delta: &total,
		})
	}

	file, err := os.CreateTemp(filepath.Dir(path), ".metrics-*.tmp")
	if err != nil {
		return fmt.Errorf("create metrics file: %w", err)
	}
	defer os.Remove(file.Name())

	if err := json.NewEncoder(file).Encode(metrics); err != nil {
		_ = file.Close()
		return fmt.Errorf("encode metrics: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close metrics file: %w", err)
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("replace metrics file: %w", err)
	}
	return nil
}

func (m *MemStorage) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read metrics: %w", err)
	}
	if len(data) == 0 {
		return nil
	}

	var metrics []models.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return fmt.Errorf("decode metrics: %w", err)
	}
	if metrics == nil {
		return fmt.Errorf("expected a JSON array of metrics")
	}

	gauges := make(map[string]float64)
	counters := make(map[string]int64)
	for i, metric := range metrics {
		if metric.ID == "" {
			return fmt.Errorf("metric at index %d has an empty id", i)
		}
		switch {
		case metric.MType == models.Gauge && metric.Value != nil && metric.Delta == nil:
			gauges[metric.ID] = *metric.Value
		case metric.MType == models.Counter && metric.Delta != nil && metric.Value == nil:
			counters[metric.ID] = *metric.Delta
		default:
			return fmt.Errorf("invalid type or value for metric %q", metric.ID)
		}
	}

	m.mu.Lock()
	m.gauges, m.counters = gauges, counters
	m.mu.Unlock()
	return nil
}

var ErrSaveSnapshot = errors.New("save metrics snapshot")

func (m *MemStorage) UpdateAndSave(path string, update func(*MemStorage) error) error {
	if err := update(m); err != nil {
		return err
	}
	if err := m.SaveToFile(path); err != nil {
		return fmt.Errorf("%w: %w", ErrSaveSnapshot, err)
	}
	return nil
}
