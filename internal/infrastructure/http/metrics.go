package http

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Metrics server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Metrics server failed: %v", err)
	}
}
