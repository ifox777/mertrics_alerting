package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifox777/metrics-alerting/internal/repository"
	"github.com/ifox777/metrics-alerting/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListHandler(t *testing.T) {
	tests := []struct {
		name            string
		gauges          map[string]float64
		counters        map[string]int64
		wantContains    []string
		wantNotContains []string
	}{
		{
			name:            "empty storage",
			wantContains:    []string{"<h1>Метрики</h1>", "<ul>"},
			wantNotContains: []string{"<li>"},
		},
		{
			name:     "both metric types",
			gauges:   map[string]float64{"temperature": 21.5},
			counters: map[string]int64{"PollCount": 5},
			wantContains: []string{
				"gauge: temperature = 21.5",
				"counter: PollCount = 5",
			},
		},
		{
			name:     "escaped metric names",
			gauges:   map[string]float64{"<script>alert(1)</script>": 21.5},
			counters: map[string]int64{"Poll&Count": 5},
			wantContains: []string{
				"gauge: &lt;script&gt;alert(1)&lt;/script&gt; = 21.5",
				"counter: Poll&amp;Count = 5",
			},
			wantNotContains: []string{"<script>", "Poll&Count"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := repository.NewMemStorage()

			for name, value := range tt.gauges {
				storage.SetGauge(name, value)
			}
			for name, value := range tt.counters {
				storage.AddCounter(name, value)
			}

			h := NewHandler(service.NewMetrics(storage))
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()

			h.ListHandler(response, request)

			require.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "text/html; charset=utf-8",
				response.Header().Get("Content-Type"))

			body := response.Body.String()

			for _, fragment := range tt.wantContains {
				assert.Contains(t, body, fragment)
			}
			for _, fragment := range tt.wantNotContains {
				assert.NotContains(t, body, fragment)
			}
		})
	}
}
