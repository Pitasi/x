package togo

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/tools/present"
)

//go:embed static templates
var embedFS embed.FS

func registerFileSlides(tree *RouteTree, ffs fs.FS, path string) error {
	tmpl := present.Template()
	if _, err := tmpl.ParseFS(embedFS, "templates/action.tmpl", "templates/slides.tmpl"); err != nil {
		panic(err)
	}

	var routePath string
	if strings.HasSuffix(path, "index.slides") {
		routePath = filepath.Dir(path)
		if routePath == "." {
			routePath = "/"
		}
	} else {
		// about.slides
		routePath = strings.TrimSuffix(path, ".slides")
	}

	tree.AddRoute(routePath, func(w http.ResponseWriter, r *http.Request) {
		f, err := ffs.Open(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		defer f.Close()

		c := &present.Context{
			ReadFile: func(filename string) ([]byte, error) {
				log.Println(filename)
				f, err := ffs.Open(filename)
				if err != nil {
					return nil, err
				}
				defer f.Close()
				return io.ReadAll(f)
			},
		}
		presentDoc, err := c.Parse(f, path, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(200)
		if err := presentDoc.Render(w, tmpl); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	})

	return nil
}
