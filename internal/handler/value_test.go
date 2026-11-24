package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ifox777/metrics-alerting/internal/repository"
	"github.com/ifox777/metrics-alerting/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestValueHandler(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantCode int
		wantBody string
	}{
		{"saved gauge", "/value/gauge/temperature", 200, "21.5"},
		{"saved zero", "/value/gauge/zero", 200, "0"},
		{"missing", "/value/gauge/missing", 404, ""},
		{"saved PollCount", "/value/counter/PollCount", 200, "5"},
		{"saved zero", "/value/counter/zero", 200, "0"},
		{"missing", "/value/counter/missing", 404, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := repository.NewMemStorage()
			storage.SetGauge("temperature", 21.5)
			storage.SetGauge("zero", 0)
			storage.AddCounter("PollCount", 2)
			storage.AddCounter("PollCount", 3)
			storage.AddCounter("zero", 0)

			h := NewHandler(service.NewMetrics(storage))

			router := chi.NewRouter()
			router.Get("/value/{type}/{name}", h.ValueHandler)
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)
			assert.Equal(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusOK {
				assert.Equal(t, tt.wantBody, response.Body.String())
				assert.Equal(t, "text/plain; charset=utf-8", response.Header().Get("Content-Type"))
			}
		})
	}
}

func TestValueJSONHandler(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.SetGauge("Alloc", 0)
	storage.SetGauge("Temperature", 21.5)
	storage.AddCounter("PollCount", 2)
	storage.AddCounter("PollCount", 3)

	h := NewHandler(service.NewMetrics(storage))
	router := chi.NewRouter()
	router.Post("/value", h.ValueJSONHandler)

	tests := []struct {
		name string
		body string
		code int
		want string
	}{
		{"zero gauge", `{"id":"Alloc","type":"gauge"}`, 200, `{"id":"Alloc","type":"gauge","value":0}`},
		{"fractional gauge", `{"id":"Temperature","type":"gauge"}`, 200, `{"id":"Temperature","type":"gauge","value":21.5}`},
		{"counter total", `{"id":"PollCount","type":"counter"}`, 200, `{"id":"PollCount","type":"counter","delta":5}`},
		{"missing metric", `{"id":"Missing","type":"gauge"}`, 404, `{"error":"metric not found"}`},
		{"unknown type", `{"id":"Alloc","type":"other"}`, 400, `{"error":"unknown metric type"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/value", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.code, resp.Code)
			assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.want, resp.Body.String())
		})
	}
}

func TestUpdateJSONHandlerCounter(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.AddCounter("PollCount", 10)
	h := NewHandler(service.NewMetrics(storage))

	body := strings.NewReader(`{"id":"PollCount","type":"counter","delta":5}`)
	req := httptest.NewRequest(http.MethodPost, "/update", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	h.UpdateJSONHandler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t,
		`{"id":"PollCount","type":"counter","delta":15}`,
		resp.Body.String(),
	)
}
