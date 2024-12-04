package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpReqs = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "HTTP requests processed, partitioned by handler.",
}, []string{"handler"})

func metricsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpReqs.WithLabelValues(r.URL.String()).Inc()
		h.ServeHTTP(w, r)
	})
}

func serveMetrics(addr string) {
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{Addr: addr, Handler: metricsMux,
		ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second}
	log.Println("Metrics server listening on", addr)
	log.Fatal(metricsServer.ListenAndServe())
}
