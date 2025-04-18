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
	"strings"
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
	ID           string
	CanonicalURL template.URL
	URL          template.URL
	Meta         ImgMeta
	Nav          ImgNav
}

func (i Img) Srcset() string {
	var sb strings.Builder
	if i.Meta.Width >= 1200 {
		sb.WriteString(string(i.CanonicalURL))
		sb.WriteString("/w_1200.webp 1200w, ")
	}
	if i.Meta.Width >= 1900 {
		sb.WriteString(string(i.CanonicalURL))
		sb.WriteString("/w_1900.webp 1900w, ")
	}
	if i.Meta.Width >= 2500 {
		sb.WriteString(string(i.CanonicalURL))
		sb.WriteString("/w_2500.webp 2500w")
	}
	return sb.String()
}

func (i Img) PreloadElement() template.HTML {
	s := fmt.Sprintf(`<link
            rel="preload" as="image" href="%s/w_1900.webp"
            imagesrcset="%s"
            imagesizes="(max-width: 1200px) 1200px, (max-width: 1900px) 1900px, 2500px">`,
		i.CanonicalURL, i.Srcset())
	return template.HTML(s)
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
	Keywords     []string
}

type ImgNav struct {
	Prev *Img
	Next *Img
}

type Images []Img

func (imgs Images) Sort() {
	sort.Slice(imgs, func(i, j int) bool { return imgs[i].Meta.Date.After(imgs[j].Meta.Date) })
}

func (imgs Images) FindByID(id string) (int, bool) {
	for i, img := range imgs {
		if img.ID == id {
			return i, true
		}
	}
	return -1, false
}

func openPhotoDB(base string) (Images, error) {
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
				ID:           id,
				Meta:         meta,
				CanonicalURL: template.URL(url),
				URL:          template.URL(url),
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

func byID(imgs []Img) map[string]int {
	m := make(map[string]int)
	for i, img := range imgs {
		m[img.ID] = i
	}
	return m
}

type ByKeywordView struct {
	keyword string
	ids     map[string]int
	imgs    []Img
}

func newByKeywordView(keyword string) *ByKeywordView {
	return &ByKeywordView{
		keyword: keyword,
		ids:     make(map[string]int),
	}
}

func (v *ByKeywordView) Append(img Img) {
	img.URL = template.URL(fmt.Sprintf("/tags/%s/pic/%s", v.keyword, img.ID))
	if len(v.imgs) > 0 {
		img.Nav.Prev = &v.imgs[len(v.imgs)-1]
		v.imgs[len(v.imgs)-1].Nav.Next = &img
	}
	v.ids[img.ID] = len(v.imgs)
	v.imgs = append(v.imgs, img)
}

func (v *ByKeywordView) Len() int {
	return len(v.imgs)
}

func (v *ByKeywordView) Images() []Img {
	return v.imgs
}

func (v *ByKeywordView) Get(id string) (Img, bool) {
	idx, ok := v.ids[id]
	if !ok {
		return Img{}, false
	}
	return v.imgs[idx], true
}

func byKeywords(imgs []Img) map[string]*ByKeywordView {
	views := make(map[string]*ByKeywordView)
	for _, img := range imgs {
		for _, keyword := range img.Meta.Keywords {
			if _, ok := views[keyword]; !ok {
				views[keyword] = newByKeywordView(keyword)
			}
			views[keyword].Append(img)
		}
	}
	return views
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

	imgs.Sort()
	imgsByID := byID(imgs)
	imgsByYear := byYear(imgs)
	imgsByKeywords := byKeywords(imgs)

	var data = Index{}
	for y, imgs := range imgsByYear {
		data.Sections = append(data.Sections, Section{
			Title: strconv.Itoa(y),
			Imgs:  imgs,
		})
	}
	sort.Slice(data.Sections, func(i, j int) bool { return data.Sections[i].Title > data.Sections[j].Title })

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
		idx, ok := imgsByID[r.PathValue("id")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		img := imgs[idx]

		if idx > 0 {
			img.Nav.Prev = &imgs[idx-1]
		}
		if idx < len(imgs)-1 {
			img.Nav.Next = &imgs[idx+1]
		}

		w.Header().Set("Cache-Control", "public, max-age=300")
		err := imagepage.Execute(w, img)
		if err != nil {
			slog.Error("rendering anto.ph image", "err", err, "id", img.ID)
		}
	})

	mux.HandleFunc("GET anto.ph/tags/{tag}", func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("tag")
		view, ok := imgsByKeywords[tag]
		if !ok {
			http.NotFound(w, r)
			return
		}

		err := homepage.Execute(w, Index{
			Sections: []Section{
				{
					Title: fmt.Sprintf("#%s (%d photos)", tag, view.Len()),
					Imgs:  view.Images(),
				},
			},
		})
		if err != nil {
			slog.Error("rendering anto.ph tag", "err", err, "tag", tag)
		}
	})

	mux.HandleFunc("GET anto.ph/tags/{tag}/pic/{id}", func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("tag")

		view, ok := imgsByKeywords[tag]
		if !ok {
			http.NotFound(w, r)
			return
		}

		img, ok := view.Get(r.PathValue("id"))
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=300")
		err := imagepage.Execute(w, img)
		if err != nil {
			slog.Error("rendering anto.ph image", "err", err, "id", img.ID)
		}
	})

	imageFilenames := []string{
		"blur.webp",
		"q_500.webp",
		"q_1000.webp",
		"q_2000.webp",
		"w_1200.webp",
		"w_1900.webp",
		"w_2500.webp",
	}
	for _, name := range imageFilenames {
		mux.HandleFunc("GET anto.ph/pic/{id}/"+name, func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			http.ServeFile(w, r, path.Join(photodbPath, id, name))
		})
	}

	mux.Handle("GET anto.ph/js/ps.js", plausible.Proxy)
	mux.Handle("GET anto.ph/api/event", plausible.Proxy)
}
