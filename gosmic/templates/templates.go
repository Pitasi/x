package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

type T struct {
	tmpl *template.Template

	// only present in devmode
	resources fs.FS
	funcs     template.FuncMap
}

func New(resources fs.FS, devmode bool, funcs template.FuncMap) *T {
	tmpl := parseFS(resources, funcs)

	if devmode {
		return &T{
			tmpl:      tmpl,
			resources: resources,
			funcs:     funcs,
		}
	}

	return &T{
		tmpl: tmpl,
	}
}

func parseFS(resources fs.FS, funcs template.FuncMap) *template.Template {
	var paths []string
	err := fs.WalkDir(resources, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(d.Name(), ".html") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return template.Must(template.New("root").Funcs(funcs).ParseFS(resources, paths...))
}

func (t *T) Reload() {
	t.tmpl = parseFS(t.resources, t.funcs)
}

func (t *T) Render(w http.ResponseWriter, name string, data any) {
	if t.resources != nil {
		t.Reload()
	}

	var buffer bytes.Buffer

	err := t.tmpl.ExecuteTemplate(&buffer, name, data)
	if err != nil {
		err = fmt.Errorf("error executing template: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_, _ = buffer.WriteTo(w)
}
