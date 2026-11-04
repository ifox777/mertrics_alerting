package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type storageSpy struct {
	gaugeName                string
	gaugeValue               float64
	counterName              string
	counterDelta             int64
	gaugeCalls, counterCalls int
}

func (s *storageSpy) SetGauge(id string, value float64) {
	s.gaugeName = id
	s.gaugeValue = value
	s.gaugeCalls++
}

func (s *storageSpy) AddCounter(id string, delta int64) {
	s.counterName = id
	s.counterDelta = delta
	s.counterCalls++
}

func (s *storageSpy) GetGauge(id string) (float64, bool) {
	return 0, false
}

func (s *storageSpy) GetCounter(id string) (int64, bool) {
	return 0, false
}

func (s *storageSpy) GetAll() (map[string]float64, map[string]int64) {
	return map[string]float64{}, map[string]int64{}
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		code        int
		wantStorage storageSpy
	}{
		{"success gauge", "POST", "/update/gauge/temperature/21.5", 200, storageSpy{gaugeName: "temperature", gaugeValue: 21.5, gaugeCalls: 1}},

		{"success counter", "POST", "/update/counter/PollCount/2", 200, storageSpy{counterName: "PollCount", counterDelta: 2, counterCalls: 1}},

		{"неправильный метод", "GET", "/update/gauge/temperature/21.5", 405, storageSpy{}},

		{"неправильный запрос 1", "POST", "/update/gauge", 404, storageSpy{}},
		{"неправильный запрос 2", "POST", "/other/gauge/temperature/1", 404, storageSpy{}},
		{"неправильный запрос 3", "POST", "/update/gauge//21.5", 404, storageSpy{}},
		{"неправильный запрос 4", "POST", "/update/unknown/temperature/1", 400, storageSpy{}},
		{"неправильный запрос 5", "POST", "/update/gauge/temperature/abc", 400, storageSpy{}},
		{"неправильный запрос 6", "POST", "/update/counter/PollCount/1.5", 400, storageSpy{}},
		{"неправильный запрос 7", "POST", "/update/gauge/temperature/", 400, storageSpy{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &storageSpy{}

			handler := NewHandler(spy)

			request := httptest.NewRequest(tt.method, tt.path, nil)
			request.Header.Set("Content-Type", "text/plain")

			response := httptest.NewRecorder()

			handler.UpdateHandler(response, request)

			assert.Equal(t, tt.code, response.Code)
			assert.Equal(t, tt.wantStorage, *spy)
		})
	}
}
