package agent

import (
	"fmt"
	"net/http"
)

func SendGauge(url string, name string, gauge float64) error {
	path := fmt.Sprintf("http://%s/update/gauge/%s/%g", url, name, gauge)

	resp, err := http.Post(path, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("send gauge %s to %s: %w", name, url, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send gauge %s to %s: unexpected HTTP status %d", name, url, resp.StatusCode)
	}

	return nil
}

func SendCounter(url string, name string, counter int64) error {
	path := fmt.Sprintf("http://%s/update/counter/%s/%d", url, name, counter)

	resp, err := http.Post(path, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("send counter %s to %s: %w", name, url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send counter %s to %s: unexpected HTTP status %d", name, url, resp.StatusCode)
	}
	return nil
}
