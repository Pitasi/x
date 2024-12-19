package main

import (
	"net/http"
	"os"
)

var host = os.Getenv("REWRITE_HOST")

func rewriteHost(h http.Handler) http.Handler {
	if len(host) == 0 {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = host
		h.ServeHTTP(w, r)
	})
}
