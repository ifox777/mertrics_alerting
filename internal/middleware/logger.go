package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)
			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("uri", r.URL.RequestURI()),
				zap.Int64("duration_us", time.Since(start).Microseconds()),
				zap.Int("status", status),
				zap.Int("bytes", ww.BytesWritten()),
			)
		})
	}
}
