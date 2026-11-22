package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/ifox777/metrics-alerting/internal/model"
)

func SendGauge(address, name string, gauge float64) error {
	return sendMetric(address, models.Metrics{
		ID:    name,
		MType: models.Gauge,
		Value: &gauge,
	})
}

func SendCounter(address, name string, delta int64) error {
	return sendMetric(address, models.Metrics{
		ID:    name,
		MType: models.Counter,
		Delta: &delta,
	})
}

func sendMetric(address string, metric models.Metrics) error {
	var body bytes.Buffer
	gz := gzip.NewWriter(&body)

	if err := json.NewEncoder(gz).Encode(metric); err != nil {
		return fmt.Errorf("encode %s %s: %w", metric.MType, metric.ID, err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("finish gzip: %w", err)
	}

	endpoint := "http://" + address + "/update"
	req, err := http.NewRequest(http.MethodPost, endpoint, &body)
	if err != nil {
		return fmt.Errorf("send %s %s to %s: %w", metric.MType, metric.ID, address, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send %s %s to %s: %w", metric.MType, metric.ID, address, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send %s %s to %s: unexpected HTTP status %d",
			metric.MType, metric.ID, address, resp.StatusCode)
	}
	return nil
}
