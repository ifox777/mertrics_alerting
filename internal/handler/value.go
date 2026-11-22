package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	models "github.com/ifox777/metrics-alerting/internal/model"
)

func (h *Handler) ValueHandler(resp http.ResponseWriter, req *http.Request) {
	value, err := h.metrics.Value(chi.URLParam(req, "type"), chi.URLParam(req, "name"))
	if err != nil {
		http.Error(resp, "Not found", http.StatusNotFound)
		return
	}

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(resp, value)
}

func (h *Handler) ValueJSONHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))

	if err != nil || mediaType != "application/json" {
		writeJSONError(resp, http.StatusUnsupportedMediaType, "expected application/json")
		return
	}
	var query models.Metrics
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&query); err != nil || query.ID == "" {
		writeJSONError(resp, http.StatusBadRequest, "expected a metric with an id")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSONError(resp, http.StatusBadRequest, "expected exactly one metric")
		return
	}
	if query.MType != models.Gauge && query.MType != models.Counter {
		writeJSONError(resp, http.StatusBadRequest, "unknown metric type")
		return
	}

	raw, err := h.metrics.Value(query.MType, query.ID)
	if err != nil {
		writeJSONError(resp, http.StatusNotFound, "metric not found")
		return
	}

	result := models.Metrics{ID: query.ID, MType: query.MType}
	switch query.MType {
	case models.Gauge:
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
			writeJSONError(resp, http.StatusInternalServerError, "invalid stored gauge")
			return
		}
		result.Value = &value
	case models.Counter:
		delta, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			writeJSONError(resp, http.StatusInternalServerError, "invalid stored counter")
			return
		}
		result.Delta = &delta
	}

	if err := json.NewEncoder(resp).Encode(result); err != nil {
		log.Printf("write metric response: %v", err)
	}
}
