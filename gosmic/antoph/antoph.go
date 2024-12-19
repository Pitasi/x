package antoph

import (
	"embed"
	"fmt"
	"gosmic/plausible"
	"html/template"
	"io/fs"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
)

//go:embed templates
var templates embed.FS

type Index struct {
	Sections []Section
}

type Section struct {
	Title string
	Imgs  []Img
}

type Img struct {
	ID   string
	URL  template.URL
	Meta ImgMeta
}

type ImgMeta struct {
	Width        int
	Height       int
	Date         time.Time
	Camera       string
	Lens         string
	ISO          string
	ShutterSpeed string
	Aperture     string
}

func openPhotoDB(base string) ([]Img, error) {
	var imgs []Img
	return imgs, fs.WalkDir(os.DirFS(base), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == "." {
			return nil
		}
		if d.IsDir() {
			id := p
			meta := extractMeta(path.Join(base, p))
			url := fmt.Sprintf("/pic/%s", id)
			imgs = append(imgs, Img{
				ID:   id,
				URL:  template.URL(url),
				Meta: meta,
			})
		}
		return nil
	})
}

func byYear(imgs []Img) map[int][]Img {
	y := make(map[int][]Img)
	for _, img := range imgs {
		y[img.Meta.Date.Year()] = append(y[img.Meta.Date.Year()], img)
	}
	return y
}

func byID(imgs []Img) map[string]Img {
	m := make(map[string]Img)
	for _, img := range imgs {
		m[img.ID] = img
	}
	return m
}

func prettyDate(t time.Time) string {
	return t.Format(time.RFC822)
}

func Register(mux *http.ServeMux) {
	photodbPath := os.Getenv("PHOTODB_PATH")
	if photodbPath == "" {
		slog.Warn("PHOTODB_PATH not set, disabling anto.ph")
		return
	}

	imgs, err := openPhotoDB(photodbPath)
	if err != nil {
		panic(err)
	}

	imgsByID := byID(imgs)

	imgsByYear := byYear(imgs)

	var data = Index{}
	for y, imgs := range imgsByYear {
		data.Sections = append(data.Sections, Section{
			Title: strconv.Itoa(y),
			Imgs:  imgs,
		})
	}
	sort.Slice(data.Sections, func(i, j int) bool { return data.Sections[i].Title > data.Sections[j].Title })
	for _, s := range data.Sections {
		sort.Slice(s.Imgs, func(i, j int) bool { return s.Imgs[i].Meta.Date.After(s.Imgs[j].Meta.Date) })
	}

	homepage := template.Must(template.ParseFS(templates, "templates/index.html"))
	imagepage, err := template.New("image.html").Funcs(template.FuncMap{
		"prettyDate": prettyDate,
	}).ParseFS(templates, "templates/image.html")
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("GET anto.ph/{$}", func(w http.ResponseWriter, r *http.Request) {
		err := homepage.Execute(w, data)
		if err != nil {
			slog.Error("rendering anto.ph homepage", "err", err)
		}
	})

	mux.HandleFunc("GET anto.ph/random/{$}", func(w http.ResponseWriter, r *http.Request) {
		random := rand.IntN(len(imgs))
		id := imgs[random].ID
		http.Redirect(w, r, fmt.Sprintf("/pic/%s", id), http.StatusFound)
	})

	mux.HandleFunc("GET anto.ph/pic/{id}", func(w http.ResponseWriter, r *http.Request) {
		img, ok := imgsByID[r.PathValue("id")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		err := imagepage.Execute(w, img)
		if err != nil {
			slog.Error("rendering anto.ph image", "err", err, "id", img.ID)
		}
	})

	mux.HandleFunc("GET anto.ph/pic/{id}/q.webp", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		http.ServeFile(w, r, path.Join(photodbPath, id, "q.webp"))
	})

	mux.HandleFunc("GET anto.ph/pic/{id}/l.webp", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		http.ServeFile(w, r, path.Join(photodbPath, id, "l.webp"))
	})

	mux.Handle("GET anto.ph/js/ps.js", plausible.Proxy)
	mux.Handle("GET anto.ph/api/event", plausible.Proxy)
}
