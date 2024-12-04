package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

type Uses struct {
	List   []App
	bySlug map[string]int
}

type App struct {
	Slug      string
	Title     string
	URL       string
	Icon      string
	Published bool
	Content   template.HTML
}

type usesFrontmatterData struct {
	Title     string
	Icon      string
	URL       string
	Published bool
}

//go:embed uses
var usesFS embed.FS

func loadUses() Uses {
	files, err := usesFS.ReadDir("uses")
	if err != nil {
		panic(err)
	}

	var list []App
	for _, file := range files {
		if file.IsDir() {
			// TODO: handle directories
			continue
		}
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		fileContent, err := usesFS.ReadFile("uses/" + file.Name())
		if err != nil {
			panic(err)
		}
		app := newUsesApp(strings.TrimSuffix(file.Name(), ".md"), fileContent)

		list = append(list, app)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Title < list[j].Title
	})

	bySlug := make(map[string]int)
	for i, app := range list {
		bySlug[app.Slug] = i
	}

	return Uses{
		List:   list,
		bySlug: bySlug,
	}
}

func newUsesApp(slug string, src []byte) App {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Linkify,
			extension.GFM,
			extension.Table,
			extension.Footnote,
			extension.Strikethrough,
			&frontmatter.Extender{},
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

	var meta usesFrontmatterData
	if err := d.Decode(&meta); err != nil {
		panic(err)
	}

	return App{
		Slug:      slug,
		Title:     meta.Title,
		URL:       meta.URL,
		Icon:      meta.Icon,
		Published: meta.Published,
		Content:   template.HTML(out.Bytes()),
	}
}

func (u Uses) Get(slug string) (App, bool) {
	idx, found := u.bySlug[slug]
	if !found {
		return App{}, false
	}
	return u.List[idx], true
}

func uses(mux *http.ServeMux) {
	uses := loadUses()

	mux.HandleFunc("GET /uses/{$}", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value("user").(*User)
		if err := templates.ExecuteTemplate(w, "uses-list.html", struct {
			Site   Site
			Active string
			User   *User
			Uses   Uses
		}{
			Site:   site,
			Active: "/uses",
			User:   user,
			Uses:   uses,
		}); err != nil {
			log.Println("executing template:", err)
		}
	})

	mux.HandleFunc("GET /uses/{slug}", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value("user").(*User)
		slug := r.PathValue("slug")
		app, found := uses.Get(slug)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err := templates.ExecuteTemplate(w, "uses.html", struct {
			Site   Site
			Active string
			User   *User
			App    App
		}{
			Site:   site,
			Active: "/uses",
			User:   user,
			App:    app,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	})

}
