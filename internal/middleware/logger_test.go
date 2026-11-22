package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLogger(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	r := chi.NewRouter()
	r.Use(RequestLogger(zap.New(core)))
	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("готово"))
	})

	req := httptest.NewRequest(http.MethodPost, "/update?source=agent", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)
	require.Equal(t, "готово", resp.Body.String())

	entries := logs.All()
	require.Len(t, entries, 1)
	require.Equal(t, zap.InfoLevel, entries[0].Level)

	fields := entries[0].ContextMap()
	require.Equal(t, "POST", fields["method"])
	require.Equal(t, "/update?source=agent", fields["uri"])
	require.EqualValues(t, resp.Code, fields["status"])
	require.EqualValues(t, len("готово"), fields["bytes"])

	duration, ok := fields["duration_us"].(int64)
	require.True(t, ok)
	require.GreaterOrEqual(t, duration, int64(0))
}
