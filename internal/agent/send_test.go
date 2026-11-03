package agent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendGauge(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{"sucess", 200, false},
		{"request error", 400, true},
		{"response error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					assert.Equal(t, r.URL.Path, "/update/gauge/Alloc/21.5")
					assert.Equal(t, r.Header.Get("Content-Type"), "text/plain")

					w.WriteHeader(tt.status)
				}))

			defer server.Close()
			err := SendGauge(strings.TrimPrefix(server.URL, "http://"), "Alloc", 21.5)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendCounter(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{"sucess", 200, false},
		{"request error", 400, true},
		{"response error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					assert.Equal(t, r.URL.Path, "/update/counter/Alloc/2")
					assert.Equal(t, r.Header.Get("Content-Type"), "text/plain")

					w.WriteHeader(tt.status)
				}))
			defer server.Close()
			err := SendCounter(strings.TrimPrefix(server.URL, "http://"), "Alloc", 2)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendGaugeNetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	))

	address := strings.TrimPrefix(server.URL, "http://")
	server.Close()

	err := SendGauge(address, "Alloc", 21.5)

	require.Error(t, err)
}

func TestSendCounterNetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	))

	address := strings.TrimPrefix(server.URL, "http://")
	server.Close()

	err := SendCounter(address, "Alloc", 2)

	require.Error(t, err)
}
