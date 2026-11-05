package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ifox777/metrics-alerting/internal/handler"
	"github.com/ifox777/metrics-alerting/internal/repository"
)

func main() {
	var server http.Server

	serverURL := flag.String("a", "localhost:8080", "input server url")
	flag.Parse()
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	r := chi.NewRouter()

	storage := repository.NewMemStorage()

	handler := handler.NewHandler(storage)

	r.Post("/update/*", handler.UpdateHandler)
	r.Get("/value/{type}/{name}", handler.ValueHandler)
	r.Get("/", handler.ListHandler)

	server.Addr = *serverURL
	server.Handler = r

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
}
