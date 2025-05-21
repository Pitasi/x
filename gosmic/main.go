package main

import (
	"html/template"
	"net/http"
	"time"

	"anto.pt/x/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"anto.pt/x/gosmic/antoph"
	"anto.pt/x/gosmic/antopt"
	"anto.pt/x/gosmic/httpx"
)

var logger = log.Module("gosmic")

func main() {
	mux := http.NewServeMux()

	httpx.RegisterWebsite("anto.pt", &antopt.Website{
		Colors: []template.CSS{
			"var(--color-default-bg)",
			"var(--color-white)",
			"var(--color-lime-200)",
			"var(--color-amber-300)",
			"var(--color-blue-200)",
			"var(--color-orange-400)",
		},
	}, mux)
	httpx.RegisterWebsite("anto.ph", antoph.Website{}, mux)

	handler := httpx.MetricsInc(mux)
	handler = httpx.RewriteHost(handler)

	s := http.Server{
		Addr:        "0.0.0.0:8080",
		Handler:     handler,
		ReadTimeout: 10 * time.Second,
	}
	go ServeMetrics(":9090")

	logger.Info("listening", "addr", s.Addr)
	logger.Error("listening", "err", s.ListenAndServe())
}

func ServeMetrics(addr string) {
	logger := log.Module("metrics_server")
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:         addr,
		Handler:      metricsMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("listening", "addr", addr)
	logger.Error("listening", "err", metricsServer.ListenAndServe())
}
