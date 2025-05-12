package httpx

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

var (
	devmode     = os.Getenv("DEVMODE") != ""
	rewriteHost = os.Getenv("REWRITE_HOST")
)

type Website interface {
	Register(devmode bool) http.Handler
}

func RegisterWebsite(domain string, website Website, mux *http.ServeMux) {
	if rewriteHost != "" && rewriteHost != domain {
		return
	}

	defer func() {
		v := recover()
		if v != nil {
			log.Printf("%s: panic: %v", domain, v)
			debug.PrintStack()
		}
	}()

	log.Println("devmode", devmode)
	h := website.Register(devmode)
	mux.Handle(domain+"/", h)
}

func RewriteHost(h http.Handler) http.Handler {
	if len(rewriteHost) == 0 {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = rewriteHost
		h.ServeHTTP(w, r)
	})
}
