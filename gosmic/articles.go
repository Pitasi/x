package main

import (
	"bytes"
	"embed"
	"fmt"
	"gosmic/md"
	"html/template"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

//go:embed articles
var articlesFS embed.FS

type Articles struct {
	list   []Article
	bySlug map[string]int
}

func loadArticles() Articles {
	var articlesFS fs.FS = articlesFS
	if devMode {
		articlesFS = os.DirFS("./")
	}

	files, err := fs.ReadDir(articlesFS, "articles")
	if err != nil {
		panic(err)
	}

	var list []Article
	for _, file := range files {
		if file.IsDir() {
			// TODO: handle directories
			continue
		}
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		fileContent, err := fs.ReadFile(articlesFS, "articles/"+file.Name())
		if err != nil {
			panic(err)
		}
		article := NewArticle(strings.TrimSuffix(file.Name(), ".md"), fileContent)

		list = append(list, article)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Date.After(list[j].Date)
	})

	bySlug := make(map[string]int)
	for i, article := range list {
		bySlug[article.Slug] = i
	}

	return Articles{
		list:   list,
		bySlug: bySlug,
	}
}

func (a Articles) Get(slug string) (Article, bool) {
	if devMode {
		a = loadArticles()
	}
	idx, found := a.bySlug[slug]
	if !found {
		return Article{}, false
	}
	return a.list[idx], true
}

func (a Articles) List() []Article {
	return a.list
}

func (a Articles) ByYear() []ArticlesByYear {
	var lastYear string
	var postsByYear []ArticlesByYear
	for _, article := range a.list {
		if !article.Published {
			continue
		}
		year := article.Date.Format("2006")
		if lastYear != year {
			lastYear = year
			postsByYear = append(postsByYear, ArticlesByYear{
				Year:  year,
				Posts: []ArticleLink{},
			})
		}
		postsByYear[len(postsByYear)-1].Posts = append(postsByYear[len(postsByYear)-1].Posts, ArticleLink{
			Title:       article.Title,
			Description: article.Description,
			Date:        article.Date.Format("02-01"),
			URL:         "/articles/" + article.Slug,
		})
	}
	return postsByYear
}

type Article struct {
	Title       string
	Slug        string
	Date        time.Time
	Description string
	Categories  []string
	Published   bool
	Content     template.HTML
}

type frontmatterData struct {
	Title       string
	Date        string
	Description string
	Categories  []string
	Published   bool
}

func NewArticle(slug string, src []byte) Article {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Linkify,
			extension.GFM,
			extension.Table,
			extension.Footnote,
			extension.Strikethrough,
			md.GosmicMarkdownExtension,
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	ctx := parser.NewContext()
	out := &bytes.Buffer{}
	md.Convert(src, out, parser.WithContext(ctx))

	d := frontmatter.Get(ctx)
	if d == nil {
		panic(fmt.Errorf("no frontmatter found in %s", slug))
	}

	var meta frontmatterData
	if err := d.Decode(&meta); err != nil {
		panic(err)
	}

	return Article{
		Slug:        slug,
		Content:     template.HTML(out.Bytes()),
		Title:       meta.Title,
		Date:        parseDate(meta.Date),
		Description: meta.Description,
		Categories:  meta.Categories,
		Published:   meta.Published,
	}
}

func parseDate(s string) time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return t
}

func articles(mux *http.ServeMux) {
	articles := loadArticles()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value("user").(*User)
		w.Header().Set("Content-Type", "text/html")
		if err := templates.ExecuteTemplate(w, "index.html", struct {
			Site     Site
			Active   string
			User     *User
			Articles []ArticlesByYear
		}{
			Site:     site,
			Active:   "/",
			User:     user,
			Articles: articles.ByYear(),
		}); err != nil {
			slog.Error("rendering template", "err", err)
		}
	})

	mux.HandleFunc("GET /articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value("user").(*User)
		slug := r.PathValue("slug")
		article, found := articles.Get(slug)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err := templates.ExecuteTemplate(w, "article.html", struct {
			Site    Site
			Active  string
			User    *User
			Article Article
		}{
			Site:    site,
			User:    user,
			Article: article,
			Active:  "/",
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	})

}
