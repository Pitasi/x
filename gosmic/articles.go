package main

import (
	"bytes"
	"embed"
	"fmt"
	"gosmic/md"
	"html/template"
	"io/fs"
	"iter"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"anto.pt/x/socialimg"
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

func (a Articles) Published() iter.Seq2[int, Article] {
	return func(yield func(int, Article) bool) {
		for i, a := range a.List() {
			if !a.Published {
				continue
			}
			if !yield(i, a) {
				return
			}
		}
	}
}

func (a Articles) ByYear() []ArticlesByYear {
	var lastYear string
	var postsByYear []ArticlesByYear
	for _, article := range a.Published() {
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

	propic, err := staticFS.Open("static/images/propic_nobg.png")
	if err != nil {
		panic(fmt.Sprintf("can't open propic_nobg.png: %s", err))
	}

	font, err := staticFS.Open("static/fonts/PPWriter-Bold.ttf")
	if err != nil {
		panic(fmt.Sprintf("can't open propic_nobg.png: %s", err))
	}

	coverGenerator, err := socialimg.NewGenerator(font, propic)
	if err != nil {
		panic(fmt.Sprintf("can't create cover generator: %s", err))
	}

	mux.HandleFunc("GET /socialimg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		if err := coverGenerator.Generate(w, "Antonio Pitasi", "https://anto.pt"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error generating social image: %v", err)
			return
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

	mux.HandleFunc("GET /articles/covers/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		article, found := articles.Get(slug)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		subtitle := fmt.Sprintf("written on %s", article.Date.Format("02 Jan 2006"))
		if err := coverGenerator.Generate(w, article.Title, subtitle); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error generating social image: %v", err)
			return
		}
	})

	feed, err := buildArticlesAtomFeed(articles)
	if err != nil {
		panic(fmt.Sprintf("can't build atom feed: %v", err))
	}

	mux.HandleFunc("GET /articles/feed.atom", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		w.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"))
		w.Write(feed)
	})
}
