package pages

import (
	"g2/antopt/uses"
	"g2/templates"
	"net/http"
)

func RenderUsesList(t *templates.T, w http.ResponseWriter, common Common, u uses.Uses) {
	type templateData struct {
		Common
		Uses uses.Uses
	}
	data := templateData{
		Common: common,
		Uses:   u,
	}
	t.Render(w, "uses-list.html", data)
}
