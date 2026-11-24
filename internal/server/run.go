package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ifox777/metrics-alerting/internal/config"
	"github.com/ifox777/metrics-alerting/internal/handler"
	"github.com/ifox777/metrics-alerting/internal/repository"
	"github.com/ifox777/metrics-alerting/internal/service"
	"go.uber.org/zap"
)

type syncMetrics struct {
	*service.Metrics
	store *repository.MemStorage
	path  string
}

func (s *syncMetrics) Update(kind, name, raw string) error {
	return s.store.UpdateAndSave(s.path, func(next *repository.MemStorage) error {
		return service.NewMetrics(next).Update(kind, name, raw)
	})
}

func Run(cfg config.ServerConfig) error {
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

	storage := repository.NewMemStorage()
	if cfg.Restore {
		if err := storage.LoadFromFile(cfg.FileStoragePath); err != nil {
			return fmt.Errorf("restore metrics from %q: %w", cfg.FileStoragePath, err)
		}
	}
	metrics := service.NewMetrics(storage)
	var api handler.MetricsService = metrics
	if cfg.StoreInterval == 0 {
		api = &syncMetrics{
			Metrics: metrics,
			store:   storage,
			path:    cfg.FileStoragePath,
		}
	}
	handler := handler.NewHandler(api)

	server.Addr = cfg.Address
	server.Handler = newRouter(handler, logger)

	serverErr := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	var ticks <-chan time.Time
	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(cfg.StoreInterval)
		defer ticker.Stop()
		ticks = ticker.C
	}

	for {
		select {
		case <-ticks:
			if err := storage.SaveToFile(cfg.FileStoragePath); err != nil {
				log.Printf("save metrics: %v", err)
			}

		case err := <-serverErr:
			return err

		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				return err
			}
			if cfg.StoreInterval > 0 {
				return storage.SaveToFile(cfg.FileStoragePath)
			}
			return nil
		}
	}
}
