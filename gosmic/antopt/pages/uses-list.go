package pages

import (
	"net/http"

	"anto.pt/x/gosmic/antopt/uses"
	"anto.pt/x/gosmic/templates"
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
