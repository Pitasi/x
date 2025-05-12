package antopt

import (
	"net/http"

	"g2/antopt/pages"
	"g2/antopt/uses"
	"g2/templates"
)

func (ws *Website) uses(t *templates.T, mux *http.ServeMux) {
	u := uses.Load()

	mux.HandleFunc("GET /uses/{$}", func(w http.ResponseWriter, r *http.Request) {
		pages.RenderUsesList(t, w, ws.common(r), u)
	})

	mux.HandleFunc("GET /uses/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		app, found := u.Get(slug)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		pages.RenderUses(t, w, ws.common(r), app)
	})
}
