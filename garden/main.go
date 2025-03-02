package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

//go:embed templates
var TemplatesFS embed.FS

//go:embed static/style.css
var StyleCSS []byte

func main() {
	obsidianVault := os.Args[1]

	files, _ := filepath.Glob(obsidianVault + "/References/*.md")

	var (
		shows  []frontmatterData
		movies []frontmatterData
		places []frontmatterData
		books  []frontmatterData
	)

	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			panic(fmt.Sprintf("reading %s: %v", f, err))
		}

		data, err := P(f, src)
		if err != nil {
			continue
		}

		if slices.Contains(data.Category, "[[Shows]]") {
			shows = append(shows, data)
		}

		if slices.Contains(data.Category, "[[Movies]]") {
			movies = append(movies, data)
		}

		if slices.Contains(data.Category, "[[Places]]") {
			places = append(places, data)
		}

		if slices.Contains(data.Category, "[[Books]]") {
			books = append(books, data)
		}
	}

	http.HandleFunc("GET /style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/css")
		w.WriteHeader(200)
		_, _ = w.Write(StyleCSS)
	})

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFS(TemplatesFS, "templates/layout.html", "templates/index.html")
		if err != nil {
			log.Println("template: index.html:", err)
			return
		}

		_ = t.Execute(w, nil)
	})

	http.HandleFunc("GET /ratings/{$}", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFS(TemplatesFS, "templates/layout.html", "templates/ratings.html")
		if err != nil {
			log.Println("template: ratings.html:", err)
			return
		}

		_ = t.Execute(w, nil)
	})

	http.HandleFunc("GET /shows/{$}", sectionHandler("shows.html", struct {
		Shows []frontmatterData
	}{
		Shows: shows,
	}))
	http.HandleFunc("GET /movies/{$}", sectionHandler("movies.html", struct {
		Movies []frontmatterData
	}{
		Movies: movies,
	}))
	http.HandleFunc("GET /places/{$}", sectionHandler("places.html", struct {
		Places []frontmatterData
	}{
		Places: places,
	}))
	http.HandleFunc("GET /books/{$}", sectionHandler("books.html", struct {
		Books []frontmatterData
	}{
		Books: books,
	}))

	log.Println("listening on", "0.0.0.0:8080")
	_ = http.ListenAndServe("0.0.0.0:8080", nil)
}

func sectionHandler(templateName string, data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("root").
			Funcs(template.FuncMap{
				"prettyDate": prettyDate,
				"yyyymmdd":   yyyymmdd,
				"prettyLink": prettyLinks,
			}).
			ParseFS(TemplatesFS, "templates/layout.html", "templates/"+templateName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unexpected error"))
			log.Println("err", err)
			return
		}

		err = t.ExecuteTemplate(w, "layout.html", data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unexpected error"))
			log.Println("err", err)
			return
		}
	}
}

type frontmatterData struct {
	Title   string
	Source  string    `yaml:"source"`
	Created time.Time `yaml:"created"`
	Tags    []string  `yaml:"tags"`

	// Movies/Shows

	Category  []string `yaml:"category"`
	Genre     []string `yaml:"genre"`
	Rating    int      `yaml:"rating"`
	ScoreIMDB float64  `yaml:"scoreImdb"`
	Cover     string   `yaml:"cover"`
	Plot      string   `yaml:"plot"`
	Year      int      `yaml:"year"`

	// Places

	Loc         []string `yaml:"loc"`
	ScoreGoogle float64  `yaml:"scoreGoogle"`
	Address     string   `yaml:"address"`
	URL         string   `yaml:"url"`

	// Books

	ScoreGR float64 `yaml:"scoreGr"`
}

func P(path string, src []byte) (frontmatterData, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Linkify,
			extension.GFM,
			extension.Table,
			extension.Footnote,
			extension.Strikethrough,
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
	if err := md.Convert(src, out, parser.WithContext(ctx)); err != nil {
		panic(err)
	}

	d := frontmatter.Get(ctx)
	if d == nil {
		return frontmatterData{}, errors.New("no frontmatter")
	}

	var meta frontmatterData
	if err := d.Decode(&meta); err != nil {
		return frontmatterData{}, err
	}

	if meta.Title == "" {
		meta.Title = pathToTitle(path)
	}

	return meta, nil
}

func prettyDate(t time.Time) string {
	if t.IsZero() {
		return "2000-01-01"
	}
	return t.Format(time.DateOnly)
}

func yyyymmdd(t time.Time) string {
	return t.Format(time.DateOnly)
}

func pathToTitle(path string) string {
	b := filepath.Base(path)
	b = strings.ReplaceAll(b, "_", " ")
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] == '.' {
			return b[0:i]
		}
	}
	return b
}

func prettyLinks(link string) string {
	return strings.TrimPrefix(strings.TrimSuffix(link, "]]"), "[[")
}
