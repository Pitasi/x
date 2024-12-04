package main

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"strings"
)

//go:embed templates
var embeddedTemplates embed.FS

type Ts struct {
	m map[string]*template.Template
}

var templates Ts

func init() {
	m := make(map[string]*template.Template)
	templatesFiles, err := fs.ReadDir(embeddedTemplates, "templates")
	if err != nil {
		panic(err)
	}
	for _, f := range templatesFiles {
		if f.IsDir() || f.Name() == "layout.html" || !strings.HasSuffix(f.Name(), ".html") {
			continue
		}
		slog.Info("registering template", "name", f.Name())
		m[f.Name()] = template.Must(template.ParseFS(embeddedTemplates, "templates/layout.html", "templates/components/*.html", "templates/"+f.Name()))
	}

	templates = Ts{
		m: m,
	}
}

func (t Ts) ExecuteTemplate(w io.Writer, name string, data any) error {
	var tmpl = t.m[name]
	if devMode {
		slog.Debug("reloading template", "name", name)
		dir := os.DirFS("./")
		tmpl = template.Must(template.ParseFS(dir, "templates/layout.html", "templates/components/*.html", "templates/"+name))
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}
