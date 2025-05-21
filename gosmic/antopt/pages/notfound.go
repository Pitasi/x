package pages

import (
	"net/http"

	"anto.pt/x/gosmic/templates"
)

func RenderNotFound(t *templates.T, w http.ResponseWriter, common Common) {
	type templateData struct {
		Common
	}
	data := templateData{
		Common: common,
	}
	t.Render(w, "notfound.html", data)
}
