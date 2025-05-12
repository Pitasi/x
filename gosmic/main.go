package main

import (
	"g2/antoph"
	"g2/antopt"
	"g2/httpx"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
		Addr:         "0.0.0.0:8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go ServeMetrics(":9090")

	log.Println("Listening on", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func ServeMetrics(addr string) {
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:         addr,
		Handler:      metricsMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Metrics server listening on", addr)
	log.Fatal(metricsServer.ListenAndServe())
}
