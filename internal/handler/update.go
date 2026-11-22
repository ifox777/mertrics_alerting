package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"

	models "github.com/ifox777/metrics-alerting/internal/model"
	"github.com/ifox777/metrics-alerting/internal/repository"
)

func (h *Handler) UpdateHandler(resp http.ResponseWriter, req *http.Request) {
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(resp, "Incorrect Path", http.StatusNotFound)
		return
	}

	if parts[1] != "update" {
		http.Error(resp, "Incorrect Path", http.StatusNotFound)
		return
	}

	if parts[3] == "" {
		http.Error(resp, "Metric name should be not empty", http.StatusNotFound)
		return
	}

	if err := h.metrics.Update(parts[2], parts[3], parts[4]); err != nil {
		if errors.Is(err, repository.ErrSaveSnapshot) {
			log.Printf("save metric: %v", err)
			http.Error(resp, "could not save metric", http.StatusInternalServerError)
			return
		}
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateJSONHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeJSONError(resp, http.StatusUnsupportedMediaType, "expected application/json")
		return
	}

	var metric models.Metrics
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&metric); err != nil || metric.ID == "" {
		writeJSONError(resp, http.StatusBadRequest, "expected a metric with an id")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSONError(resp, http.StatusBadRequest, "expected exaxtly one metric")
		return
	}

	var rawValue string
	switch {
	case metric.MType == models.Gauge && metric.Value != nil:
		rawValue = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case metric.MType == models.Counter && metric.Delta != nil:
		rawValue = strconv.FormatInt(*metric.Delta, 10)
	default:
		writeJSONError(resp, http.StatusBadRequest, "expected gauge with value or counter with delta")
		return
	}

	if err := h.metrics.Update(metric.MType, metric.ID, rawValue); err != nil {
		log.Printf("update metric %q: %v", metric.ID, err)
		writeJSONError(resp, http.StatusInternalServerError, "could not update metric")
		return
	}
	if err := json.NewEncoder(resp).Encode(metric); err != nil {
		log.Printf("write metric response: %v", err)
	}
}

func writeJSONError(resp http.ResponseWriter, status int, message string) {
	resp.WriteHeader(status)
	if err := json.NewEncoder(resp).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("write JSON error: %v", err)
	}
}
