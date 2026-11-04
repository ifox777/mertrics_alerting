package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ValueHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(resp, "Incorrect method", http.StatusMethodNotAllowed)
		return
	}

	switch chi.URLParam(req, "type") {
	case "gauge":
		id := chi.URLParam(req, "name")
		value, has := h.storage.GetGauge(id)
		if !has {
			http.Error(resp, "Not found", http.StatusNotFound)
			return
		}
		resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(resp, value)
	case "counter":
		id := chi.URLParam(req, "name")
		value, has := h.storage.GetCounter(id)
		if !has {
			http.Error(resp, "Not found", http.StatusNotFound)
			return
		}
		resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(resp, value)

	default:
		http.Error(resp, "Not found", http.StatusNotFound)
		return

	}
}
