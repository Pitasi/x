package httpx

import (
	"net/http"
	"os"
	"runtime/debug"

	"anto.pt/x/log"
)

var (
	devmode     = os.Getenv("DEVMODE") != ""
	rewriteHost = os.Getenv("REWRITE_HOST")
)

var logger = log.Module("httpx")

type Website interface {
	Register(devmode bool) http.Handler
}

func RegisterWebsite(domain string, website Website, mux *http.ServeMux) {
	if rewriteHost != "" && rewriteHost != domain {
		return
	}

	logger := logger.With("domain", domain)

	defer func() {
		v := recover()
		if v != nil {
			logger.Error("panic recovered", "err", v, "stack", debug.Stack())
		}
	}()

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
