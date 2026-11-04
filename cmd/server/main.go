package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ifox777/metrics-alerting/internal/handler"
	"github.com/ifox777/metrics-alerting/internal/repository"
)

func main() {
	r := chi.NewRouter()

	storage := repository.NewMemStorage()

	handler := handler.NewHandler(storage)

	r.Post("/update/*", handler.UpdateHandler)
	r.Get("/value/{type}/{name}", handler.ValueHandler)
	r.Get("/", handler.ListHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
