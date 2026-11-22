package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ifox777/metrics-alerting/internal/config"
	"github.com/ifox777/metrics-alerting/internal/handler"
	"github.com/ifox777/metrics-alerting/internal/middleware"
	"github.com/ifox777/metrics-alerting/internal/repository"
	"github.com/ifox777/metrics-alerting/internal/service"
	"go.uber.org/zap"
)

func run(cfg config.ServerConfig) error {
	var server http.Server
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	logConfig := zap.NewProductionConfig()
	logConfig.Sampling = nil

	logger, err := logConfig.Build()
	if err != nil {
		return fmt.Errorf("create HTTP logger: %w", err)
	}
	defer func() { _ = logger.Sync() }()
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger))

	storage := repository.NewMemStorage()
	metrics := service.NewMetrics(storage)
	handler := handler.NewHandler(metrics)

	r.Post("/update/*", handler.UpdateHandler)
	r.Post("/update", handler.UpdateJSONHandler)
	r.Post("/value", handler.ValueJSONHandler)
	r.Post("/update/", handler.UpdateJSONHandler)
	r.Post("/value/", handler.ValueJSONHandler)
	r.Get("/value/{type}/{name}", handler.ValueHandler)
	r.Get("/", handler.ListHandler)

	server.Addr = cfg.Address
	server.Handler = r

	serverErr := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		if err == http.ErrServerClosed {
			return nil
		}
		return err

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
			return err
		}
		return nil
	}
}

func main() {
	cfg, err := config.ParseServer(os.Args[1:])
	if err == flag.ErrHelp {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}
