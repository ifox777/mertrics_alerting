package main

import (
	"flag"
	"log"
	"os"

	"github.com/ifox777/metrics-alerting/internal/config"
	"github.com/ifox777/metrics-alerting/internal/server"
)

func main() {
	cfg, err := config.ParseServer(os.Args[1:])
	if err == flag.ErrHelp {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
