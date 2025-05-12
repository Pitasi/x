package uses

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
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

//go:embed *.md
var usesFS embed.FS

func Load() Uses {
	files, err := usesFS.ReadDir(".")
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

		fileContent, err := usesFS.ReadFile(file.Name())
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
	if err := md.Convert(src, out, parser.WithContext(ctx)); err != nil {
		panic(err)
	}

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
