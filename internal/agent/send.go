package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/ifox777/metrics-alerting/internal/model"
)

func SendGauge(address, name string, gauge float64) error {
	metric := models.Metrics{
		ID:    name,
		MType: models.Gauge,
		Value: &gauge,
	}

	data, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("encode gauge %s: %w", name, err)
	}

	endpoint := "http://" + address + "/update"
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("send gauge %s to %s: %w", name, address, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send gauge %s to %s: unexpected HTTP status %d", name, address, resp.StatusCode)
	}
	return nil
}

func SendCounter(address, name string, delta int64) error {
	metric := models.Metrics{
		ID:    name,
		MType: models.Counter,
		Delta: &delta,
	}

	data, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("encode counter %s: %w", name, err)
	}

	endpoint := "http://" + address + "/update"
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("send counter %s to %s: %w", name, address, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send counter %s to %s: unexpected HTTP status %d", name, address, resp.StatusCode)
	}
	return nil
}
