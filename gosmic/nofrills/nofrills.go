package nofrills

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates
var templates embed.FS

func Register(mux *http.ServeMux) {
	homepage := template.Must(template.ParseFS(templates, "templates/index.html"))

	mux.HandleFunc("GET nofrills.systems/", func(w http.ResponseWriter, r *http.Request) {
		homepage.Execute(w, nil)
	})
}
