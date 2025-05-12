package antopt

import (
	p "g2/plausible"
	"net/http"
)

func (ws *Website) plausible(mux *http.ServeMux) {
	mux.Handle("/js/ps.js", p.Proxy)
	mux.Handle("/api/event", p.Proxy)
}
