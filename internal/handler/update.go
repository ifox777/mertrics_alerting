package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) UpdateHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(resp, "Incorrect method", http.StatusMethodNotAllowed)
		return
	}
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

	switch parts[2] {
	case "gauge":
		gauge, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			http.Error(resp, "Ivalid gauge", http.StatusBadRequest)
			return
		}
		h.storage.SetGauge(parts[3], gauge)
	case "counter":
		counter, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			http.Error(resp, "Invalid counter", http.StatusBadRequest)
			return
		}
		h.storage.AddCounter(parts[3], counter)
	default:
		http.Error(resp, "Incorrect metric", http.StatusBadRequest)
		return
	}

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
}
