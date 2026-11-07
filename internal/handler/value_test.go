package handler

import (
	"net/http"
	"net/http/httptest"
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
