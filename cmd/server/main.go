package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ifox777/metrics-alerting/internal/handler"
	"github.com/ifox777/metrics-alerting/internal/repository"
)

func main() {
	serverURL := flag.String("a", "localhost:8080", "input server url")
	flag.Parse()
	r := chi.NewRouter()

	storage := repository.NewMemStorage()

	handler := handler.NewHandler(storage)

	r.Post("/update/*", handler.UpdateHandler)
	r.Get("/value/{type}/{name}", handler.ValueHandler)
	r.Get("/", handler.ListHandler)

	err := http.ListenAndServe(*serverURL, r)
	if err != nil {
		panic(err)
	}
}
