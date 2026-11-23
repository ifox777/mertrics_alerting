package handler

import (
	"net/http"
	"strings"
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
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
}
