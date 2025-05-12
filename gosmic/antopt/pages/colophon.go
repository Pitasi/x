package pages

import (
	"g2/templates"
	"net/http"
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
