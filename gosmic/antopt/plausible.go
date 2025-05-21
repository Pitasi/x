package antopt

import (
	"net/http"

	p "anto.pt/x/gosmic/plausible"
)

func (ws *Website) plausible(mux *http.ServeMux) {
	mux.Handle("/js/ps.js", p.Proxy)
	mux.Handle("/api/event", p.Proxy)
}
