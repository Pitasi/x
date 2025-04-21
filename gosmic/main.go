package main

import (
	"log"
	"net/http"
	"os"

	"gosmic/antoph"
	"gosmic/lastfm"
	"gosmic/nofrills"
)

var (
	devMode      = os.Getenv("DEV_MODE") == "true"
	lastfmAPIKey = os.Getenv("LASTFM_API_KEY")
)

func main() {
	addr := "0.0.0.0:8080"
	go serveMetrics(":9090")

	mux := http.NewServeMux()

	if lastfmAPIKey != "" {
		c := &lastfm.Client{
			APIKey: lastfmAPIKey,
		}
		NowPlaying(c, mux)
	}

	articles(mux)
	uses(mux)
	colophon(mux)
	x(mux)
	redirects(mux)
	prefs(mux)
	static(mux)
	plausible(mux)
	notfound(mux)
	nofrills.Register(mux)
	antoph.Register(mux)

	var h http.Handler = mux
	h = recoverer(h)
	h = logger(h)
	h = compress(h)
	h = userMiddleware(h)
	h = metricsMiddleware(h)
	if devMode {
		h = rewriteHost(h)
	}

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, h))
}
