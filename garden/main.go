package main

import (
	"bytes"
	_ "embed"
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

//go:embed index.html
var IndexHtml string

//go:embed movies.html
var MoviesHtml string

//go:embed shows.html
var ShowsHtml string

func main() {
	obsidianVault := os.Args[1]

	files, _ := filepath.Glob(obsidianVault + "/References/*.md")

	var (
		shows  []frontmatterData
		movies []frontmatterData
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
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index").
			Parse(string(IndexHtml))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unexpected error"))
			log.Println("err", err)
			return
		}

		err = t.Execute(w, struct{}{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unexpected error"))
			log.Println("err", err)
			return
		}
	})

	http.HandleFunc("GET /shows", sectionHandler(ShowsHtml, struct {
		Shows []frontmatterData
	}{
		Shows: shows,
	}))
	http.HandleFunc("GET /movies", sectionHandler(MoviesHtml, struct {
		Movies []frontmatterData
	}{
		Movies: movies,
	}))

	log.Println("listening on", "0.0.0.0:8080")
	_ = http.ListenAndServe("0.0.0.0:8080", nil)
}

func sectionHandler(text string, data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index").
			Funcs(template.FuncMap{
				"prettyDate": prettyDate,
				"yyyymmdd":   yyyymmdd,
			}).
			Parse(text)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unexpected error"))
			log.Println("err", err)
			return
		}

		err = t.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unexpected error"))
			log.Println("err", err)
			return
		}
	}
}

type frontmatterData struct {
	Title     string
	Category  []string  `yaml:"category"`
	Genre     []string  `yaml:"genre"`
	Rating    int       `yaml:"rating"`
	ScoreIMDB float64   `yaml:"scoreImdb"`
	Cover     string    `yaml:"cover"`
	Plot      string    `yaml:"plot"`
	Year      int       `yaml:"year"`
	Created   time.Time `yaml:"created"`
	Tags      []string  `yaml:"tags"`
	Source    string    `yaml:"source"`
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

	meta.Title = pathToTitle(path)
	return meta, nil
}

func prettyDate(t time.Time) string {
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
