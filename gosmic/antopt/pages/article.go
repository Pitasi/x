package pages

import (
	"net/http"

	"anto.pt/x/gosmic/antopt/articles"
	"anto.pt/x/gosmic/templates"
)

func RenderArticle(t *templates.T, w http.ResponseWriter, common Common, article articles.Article) {
	type templateData struct {
		Common

		Article articles.Article
	}

	tmplData := templateData{
		Common:  common,
		Article: article,
	}

	t.Render(w, "article.html", tmplData)
}
