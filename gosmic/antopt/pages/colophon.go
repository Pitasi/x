package pages

import (
	"net/http"

	"anto.pt/x/gosmic/templates"
)

func RenderColophon(t *templates.T, w http.ResponseWriter, common Common) {
	type templateData struct {
		Common
	}
	data := templateData{
		Common: common,
	}
	t.Render(w, "colophon.html", data)
}
