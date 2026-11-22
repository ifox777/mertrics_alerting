package config

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Address         string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
}

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func ParseServer(args []string) (ServerConfig, error) {
	fs := flag.NewFlagSet("flag", flag.ContinueOnError)
	address := fs.String("a", "localhost:8080", "server address")
	interval := fs.String("i", "300", "store interval in seconds")
	filePath := fs.String("f", "metrics.json", "metrics storage file")
	restore := fs.Bool("r", true, "restore metrics from file on startup")

	if err := fs.Parse(args); err != nil {
		return ServerConfig{}, err
	}

	rawInterval := *interval
	if value, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		rawInterval = value
	}

	seconds, err := strconv.ParseInt(rawInterval, 10, 64)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("parse store interval %q: %w", rawInterval, err)
	}
	const maxSeconds = math.MaxInt64 / int64(time.Second)
	if seconds < 0 || seconds > maxSeconds {
		return ServerConfig{}, fmt.Errorf("store interval must be between 0 and %d seconds", maxSeconds)
	}

	cfg := ServerConfig{
		Address:         *address,
		StoreInterval:   time.Duration(seconds) * time.Second,
		FileStoragePath: *filePath,
		Restore:         *restore,
	}
	if value, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Address = value
	}
	if value, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = value
	}
	if raw, ok := os.LookupEnv("RESTORE"); ok {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return ServerConfig{}, fmt.Errorf("parse RESTORE %q: %w", raw, err)
		}
		cfg.Restore = value
	}
	return cfg, nil
}

func ParseAgent(args []string) (AgentConfig, error) {
	fs := flag.NewFlagSet("flag", flag.ContinueOnError)

	address := fs.String("a", "localhost:8080", "server address")
	reportInterval := fs.Int("r", 10, "report interval")
	pollinterval := fs.Int("p", 2, "poll interval")

	err := fs.Parse(args)
	if err != nil {
		return AgentConfig{}, err
	}

	cfg := AgentConfig{
		Address:        *address,
		ReportInterval: time.Duration(*reportInterval) * time.Second,
		PollInterval:   time.Duration(*pollinterval) * time.Second,
	}

	envAddress, ok := os.LookupEnv("ADDRESS")
	if ok {
		cfg.Address = envAddress
	}

	envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		seconds, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return AgentConfig{}, fmt.Errorf("parse REPORT_INTERVAL: %w", err)
		}
		cfg.ReportInterval = time.Duration(seconds) * time.Second
	}

	envPollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		seconds, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return AgentConfig{}, fmt.Errorf("parse POLL_INTERVAL: %w", err)
		}
		cfg.PollInterval = time.Duration(seconds) * time.Second

	}

	return cfg, nil
}
