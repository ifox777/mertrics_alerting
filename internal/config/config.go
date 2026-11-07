package config

import (
	"flag"
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

	return ServerConfig{
		Address: *address,
	}, nil
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
	return AgentConfig{
		Address:        *address,
		ReportInterval: time.Duration(*reportInterval) * time.Second,
		PollInterval:   time.Duration(*pollinterval) * time.Second,
	}, nil
}
