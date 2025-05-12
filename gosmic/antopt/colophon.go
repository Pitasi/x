package antopt

import (
	"g2/antopt/pages"
	"g2/templates"
	"net/http"
)

func (ws *Website) colophon(t *templates.T, mux *http.ServeMux) {
	mux.HandleFunc("GET /colophon/{$}", func(w http.ResponseWriter, r *http.Request) {
		pages.RenderColophon(t, w, ws.common(r))
	})
}
