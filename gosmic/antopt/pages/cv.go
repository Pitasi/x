package pages

import (
	"net/http"

	"anto.pt/x/gosmic/templates"
)

func RenderCV(t *templates.T, w http.ResponseWriter, common Common) {
	type templateData struct {
		Common
	}
	data := templateData{
		Common: common,
	}
	t.Render(w, "cv.html", data)
}
