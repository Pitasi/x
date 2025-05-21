package pages

import (
	"net/http"

	"anto.pt/x/gosmic/antopt/uses"
	"anto.pt/x/gosmic/templates"
)

func RenderUses(t *templates.T, w http.ResponseWriter, common Common, app uses.App) {
	type templateData struct {
		Common
		App uses.App
	}
	data := templateData{
		Common: common,
		App:    app,
	}
	t.Render(w, "uses.html", data)
}
