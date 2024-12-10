package main

import (
	"net/http"

	p "gosmic/plausible"
)

func plausible(mux *http.ServeMux) {
	mux.Handle("/js/ps.js", p.Proxy)
	mux.Handle("/api/event", p.Proxy)
}
