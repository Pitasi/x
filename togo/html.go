package togo

import (
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

func registerFileHTML(tree *RouteTree, ffs fs.FS, path string) error {
	f, err := ffs.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
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

	if _, err := tml.Parse(string(content)); err != nil {
		return err
	}

	var routePath string
	if strings.HasSuffix(path, "index.html") {
		routePath = filepath.Dir(path)
		if routePath == "." {
			routePath = "/"
		}
	} else {
		// e.g. about.html
		routePath = strings.TrimSuffix(path, ".html")
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
