package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ifox777/metrics-alerting/internal/handler"
	"github.com/ifox777/metrics-alerting/internal/middleware"
	"go.uber.org/zap"
)

func newRouter(h *handler.Handler, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(logger))
	r.Use(chimiddleware.Compress(5, "application/json", "text/html"))
	r.Use(middleware.DecompressRequest)

	r.Post("/update/*", h.UpdateHandler)
	r.Post("/update", h.UpdateJSONHandler)
	r.Post("/value", h.ValueJSONHandler)
	r.Post("/update/", h.UpdateJSONHandler)
	r.Post("/value/", h.ValueJSONHandler)
	r.Get("/value/{type}/{name}", h.ValueHandler)
	r.Get("/", h.ListHandler)

	return r
}
