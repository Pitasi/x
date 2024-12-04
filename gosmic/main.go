package main

import (
	"log"
	"net/http"
	"os"
)

var devMode = os.Getenv("DEV_MODE") == "true"

func main() {
	addr := "0.0.0.0:8080"
	go serveMetrics(":9090")

	mux := http.NewServeMux()

	articles(mux)
	uses(mux)
	x(mux)
	redirects(mux)
	prefs(mux)
	static(mux)
	plausible(mux)
	notfound(mux)

	var h http.Handler = mux
	h = recoverer(h)
	h = logger(h)
	h = compress(h)
	h = userMiddleware(h)
	h = metricsMiddleware(h)

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, h))
}
