package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Address string
}

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func ParseServer(args []string) (ServerConfig, error) {
	fs := flag.NewFlagSet("flag", flag.ContinueOnError)

	address := fs.String("a", "localhost:8080", "server address")

	err := fs.Parse(args)
	if err != nil {
		return ServerConfig{}, err
	}

	cfg := ServerConfig{Address: *address}

	envAddress, ok := os.LookupEnv("ADDRESS")
	if ok {
		cfg.Address = envAddress
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
