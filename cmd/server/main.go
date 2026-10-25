package main

import (
	"net/http"

	"github.com/ifox777/metrics-alerting/internal/handler"
	"github.com/ifox777/metrics-alerting/internal/repository"
)

func main() {
	storage := repository.NewMemStorage()

	handler := handler.NewHandler(storage)

	err := http.ListenAndServe(":8080", http.HandlerFunc(handler.UpdateHandler))
	if err != nil {
		panic(err)
	}
}
