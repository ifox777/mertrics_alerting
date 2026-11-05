package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"runtime"
	"time"

	"github.com/ifox777/metrics-alerting/internal/agent"
)

func main() {
	var metric runtime.MemStats
	var pollCount int64

	agentURL := flag.String("a", "localhost:8080", "input agent url")
	reportInterval := flag.Int("r", 10, "report interval")
	pollInterval := flag.Int("p", 2, "poll interval")

	flag.Parse()

	lastReport := time.Now()

	var lastSentPollCount int64

	var metrics map[string]float64
	var lastPoll time.Time

	for {
		if lastPoll.IsZero() || time.Since(lastPoll) >= time.Duration(*pollInterval)*time.Second {
			runtime.ReadMemStats(&metric)

			metrics = map[string]float64{
				"Alloc":         float64(metric.Alloc),
				"BuckHashSys":   float64(metric.BuckHashSys),
				"Frees":         float64(metric.Frees),
				"GCCPUFraction": float64(metric.GCCPUFraction),
				"GCSys":         float64(metric.GCSys),
				"HeapAlloc":     float64(metric.HeapAlloc),
				"HeapIdle":      float64(metric.HeapIdle),
				"HeapInuse":     float64(metric.HeapInuse),
				"HeapObjects":   float64(metric.HeapObjects),
				"HeapReleased":  float64(metric.HeapReleased),
				"HeapSys":       float64(metric.HeapSys),
				"LastGC":        float64(metric.LastGC),
				"Lookups":       float64(metric.Lookups),
				"MCacheInuse":   float64(metric.MCacheInuse),
				"MCacheSys":     float64(metric.MCacheSys),
				"MSpanInuse":    float64(metric.MSpanInuse),
				"MSpanSys":      float64(metric.MSpanSys),
				"Mallocs":       float64(metric.Mallocs),
				"NextGC":        float64(metric.NextGC),
				"NumForcedGC":   float64(metric.NumForcedGC),
				"NumGC":         float64(metric.NumGC),
				"OtherSys":      float64(metric.OtherSys),
				"PauseTotalNs":  float64(metric.PauseTotalNs),
				"StackInuse":    float64(metric.StackInuse),
				"StackSys":      float64(metric.StackSys),
				"TotalAlloc":    float64(metric.TotalAlloc),
				"Sys":           float64(metric.Sys),
				"RandomValue":   rand.Float64(),
			}

			pollCount = pollCount + 1
			lastPoll = time.Now()
		}
		if time.Since(lastReport) >= time.Duration(*reportInterval)*time.Second {
			for name, metric := range metrics {
				err := agent.SendGauge(*agentURL, name, metric)
				if err != nil {
					log.Println(err)
				}
			}

			err := agent.SendCounter(*agentURL, "PollCount", pollCount-lastSentPollCount)
			if err != nil {
				log.Println(err)
			} else {
				lastSentPollCount = pollCount
			}
			lastReport = time.Now()

		}

		fmt.Println(metrics)
		fmt.Println(pollCount)

		time.Sleep(time.Second)
	}
}
