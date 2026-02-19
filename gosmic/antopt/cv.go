package antopt

import (
	"net/http"

	"anto.pt/x/gosmic/antopt/pages"
	"anto.pt/x/gosmic/templates"
)

func (ws *Website) cv(t *templates.T, mux *http.ServeMux) {
	mux.HandleFunc("GET /cv/{$}", func(w http.ResponseWriter, r *http.Request) {
		pages.RenderCV(t, w, ws.common(r))
	})
}
