package togo

import (
	"bytes"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

func registerFileMD(tree *RouteTree, ffs fs.FS, path string) error {
	f, err := ffs.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	mdContent, err := parseMarkdown(content)
	if err != nil {
		return err
	}

	tml := template.New(":root").Funcs(template.FuncMap{
		"PathParam": func(_ string) string {
			return "oof this is no good"
		},
	})

	var layouts []string
	layoutPath := filepath.Join(filepath.Dir(path), "_layout.html")
	for {
		if exists(ffs, layoutPath) {
			layoutF, err := ffs.Open(layoutPath)
			if err != nil {
				return err
			}
			layoutContent, err := io.ReadAll(layoutF)
			if err != nil {
				return err
			}

			layouts = append(layouts, string(layoutContent))
		}

		if layoutPath == filepath.Join(".", "_layout.html") {
			break
		}

		parent := filepath.Dir(filepath.Dir(layoutPath))
		layoutPath = filepath.Join(parent, "_layout.html")
	}

	slices.Reverse(layouts)
	for _, c := range layouts {
		_, err = tml.Parse(c)
		if err != nil {
			return err
		}
	}

	mdTml, err := template.New("body").Parse(string(mdContent))
	if err != nil {
		return err
	}

	_, err = tml.AddParseTree("body", mdTml.Tree)
	if err != nil {
		return err
	}

	var routePath string
	if strings.HasSuffix(path, "index.md") {
		routePath = filepath.Dir(path)
		if routePath == "." {
			routePath = "/"
		}
	} else {
		// about.md
		routePath = strings.TrimSuffix(path, ".md")
	}

	tree.AddRoute(routePath, func(w http.ResponseWriter, r *http.Request) {
		if err := tml.Funcs(template.FuncMap{
			"PathParam": func(name string) string {
				return r.PathValue(name)
			},
		}).Execute(w, struct {
			Globals map[string]any
		}{
			Globals: map[string]any{
				"var": 42,
			},
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	})

	return nil
}

func parseMarkdown(src []byte) (template.HTML, error) {
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
		return "", err
	}

	return template.HTML(out.Bytes()), nil
}
