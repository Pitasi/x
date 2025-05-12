package pages

import (
	"g2/antopt/uses"
	"g2/templates"
	"net/http"
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
