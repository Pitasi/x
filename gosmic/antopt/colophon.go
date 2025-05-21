package antopt

import (
	"net/http"

	"anto.pt/x/gosmic/antopt/pages"
	"anto.pt/x/gosmic/templates"
)

func (ws *Website) colophon(t *templates.T, mux *http.ServeMux) {
	mux.HandleFunc("GET /colophon/{$}", func(w http.ResponseWriter, r *http.Request) {
		pages.RenderColophon(t, w, ws.common(r))
	})
}
