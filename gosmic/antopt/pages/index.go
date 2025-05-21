package pages

import (
	"net/http"

	"anto.pt/x/gosmic/antopt/articles"
	"anto.pt/x/gosmic/templates"
)

func RenderIndex(t *templates.T, w http.ResponseWriter, common Common, a []articles.ArticlesByYear) {
	type templateData struct {
		Common

		CurrentURL string
		Articles   []articles.ArticlesByYear
	}

	tmplData := templateData{
		Common:   common,
		Articles: a,
	}

	t.Render(w, "index.html", tmplData)
}
