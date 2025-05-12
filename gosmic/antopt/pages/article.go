package pages

import (
	"g2/antopt/articles"
	"g2/templates"
	"net/http"
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
