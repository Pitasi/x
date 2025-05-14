package antoph

import (
	"embed"
	"fmt"
	"g2/fsx"
	"g2/httpx"
	"g2/plausible"
	"g2/templates"
	"html/template"
	"io/fs"
	"log"
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

//go:embed pages/*.html
var resources embed.FS

var featured []string

func init() {
	for s := range strings.SplitSeq(os.Getenv("FEATURED"), ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			featured = append(featured, s)
		}
	}
}

type Index struct {
	Sections []Section
	Tags     []string
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
			if keyword == "" {
				log.Printf("img %s has empty keyword", img.ID)
				continue
			}
			if _, ok := views[keyword]; !ok {
				views[keyword] = newByKeywordView(keyword)
			}
			views[keyword].Append(img)
		}
	}
	return views
}

type AllView struct {
	ids  map[string]int
	imgs []Img
}

func (v *AllView) Append(img Img) {
	img.URL = template.URL(fmt.Sprintf("/all/pic/%s", img.ID))
	if len(v.imgs) > 0 {
		img.Nav.Prev = &v.imgs[len(v.imgs)-1]
		v.imgs[len(v.imgs)-1].Nav.Next = &img
	}
	v.ids[img.ID] = len(v.imgs)
	v.imgs = append(v.imgs, img)
}

func (v *AllView) Len() int {
	return len(v.imgs)
}

func (v *AllView) Images() []Img {
	return v.imgs
}

func (v *AllView) Get(id string) (Img, bool) {
	idx, ok := v.ids[id]
	if !ok {
		return Img{}, false
	}
	return v.imgs[idx], true
}

func (v *AllView) ByYear() Index {
	var data = Index{}
	for y, imgs := range byYear(v.imgs) {
		data.Sections = append(data.Sections, Section{
			Title: strconv.Itoa(y),
			Imgs:  imgs,
		})
	}
	sort.Slice(data.Sections, func(i, j int) bool { return data.Sections[i].Title > data.Sections[j].Title })
	return data
}

func newAllView(imgs []Img) AllView {
	v := AllView{
		ids:  make(map[string]int),
		imgs: make([]Img, 0, len(imgs)),
	}
	for _, i := range imgs {
		v.Append(i)
	}
	return v
}

type FeaturedView struct {
	ids  map[string]int
	imgs []Img
}

func newFeaturedView(imgs []Img) *FeaturedView {
	v := &FeaturedView{
		ids:  make(map[string]int),
		imgs: make([]Img, 0, len(imgs)),
	}

	for _, img := range imgs {
		for _, feat := range featured {
			if img.ID == feat {
				v.Append(img)
			}
		}

	}

	return v
}

func (v *FeaturedView) Images() []Img { return v.imgs }

func (v *FeaturedView) Get(id string) (Img, bool) {
	idx, ok := v.ids[id]
	if !ok {
		return Img{}, false
	}
	return v.imgs[idx], true
}

func (v *FeaturedView) Append(img Img) {
	img.URL = template.URL(fmt.Sprintf("/pic/%s", img.ID))
	if len(v.imgs) > 0 {
		img.Nav.Prev = &v.imgs[len(v.imgs)-1]
		v.imgs[len(v.imgs)-1].Nav.Next = &img
	}
	v.ids[img.ID] = len(v.imgs)
	v.imgs = append(v.imgs, img)
}

func prettyDate(t time.Time) string {
	return t.Format(time.RFC822)
}

type Website struct{}

var _ httpx.Website = Website{}

func (Website) Register(devmode bool) http.Handler {
	mux := http.NewServeMux()

	photodbPath := os.Getenv("PHOTODB_PATH")
	if photodbPath == "" {
		slog.Warn("PHOTODB_PATH not set")
		return nil
	}

	imgs, err := openPhotoDB(photodbPath)
	if err != nil {
		panic(err)
	}

	imgs.Sort()
	imgsByYear := byYear(imgs)
	imgsByKeywords := byKeywords(imgs)

	var keywords []string
	for k := range imgsByKeywords {
		keywords = append(keywords, k)
	}
	sort.Strings(keywords)

	var data = Index{
		Tags: keywords,
	}
	for y, imgs := range imgsByYear {
		data.Sections = append(data.Sections, Section{
			Title: strconv.Itoa(y),
			Imgs:  imgs,
		})
	}
	sort.Slice(data.Sections, func(i, j int) bool { return data.Sections[i].Title > data.Sections[j].Title })

	t := templates.New(fsx.Or(devmode, resources, "./antoph"), devmode, template.FuncMap{
		"prettyDate": prettyDate,
	})

	featuredView := newFeaturedView(imgs)

	if len(featured) == 0 {
		// fallback: homepage show all the images
		mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
			t.Render(w, "index.html", data)
		})
	} else {
		mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
			t.Render(w, "index.html", Index{
				Sections: []Section{{Title: "Top pic(k)s", Imgs: featuredView.Images()}},
				Tags:     keywords,
			})
		})
	}

	allView := newAllView(imgs)
	allViewData := allView.ByYear()

	mux.HandleFunc("GET /all/{$}", func(w http.ResponseWriter, r *http.Request) {
		t.Render(w, "index.html", allViewData)
	})

	mux.HandleFunc("GET /all/pic/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		img, ok := allView.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}

		if !devmode {
			w.Header().Set("Cache-Control", "public, max-age=300")
		}
		t.Render(w, "image.html", img)
	})

	mux.HandleFunc("GET /random/{$}", func(w http.ResponseWriter, r *http.Request) {
		random := rand.IntN(len(imgs))
		id := imgs[random].ID
		http.Redirect(w, r, fmt.Sprintf("/all/pic/%s", id), http.StatusFound)
	})

	mux.HandleFunc("GET /pic/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		img, ok := featuredView.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}

		if !devmode {
			w.Header().Set("Cache-Control", "public, max-age=300")
		}
		t.Render(w, "image.html", img)
	})

	mux.HandleFunc("GET /tags/{tag}/", func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("tag")
		view, ok := imgsByKeywords[tag]
		if !ok {
			http.NotFound(w, r)
			return
		}

		t.Render(w, "index.html", Index{
			Sections: []Section{
				{
					Title: fmt.Sprintf("#%s (%d photos)", tag, view.Len()),
					Imgs:  view.Images(),
				},
			},
		})
	})

	mux.HandleFunc("GET /tags/{tag}/pic/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		if !devmode {
			w.Header().Set("Cache-Control", "public, max-age=300")
		}
		t.Render(w, "image.html", img)
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
		mux.HandleFunc("GET /pic/{id}/"+name, func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			if !devmode {
				w.Header().Set("Cache-Control", "public, max-age=31536000")
			}
			http.ServeFile(w, r, path.Join(photodbPath, id, name))
		})
	}

	mux.Handle("GET /js/ps.js", plausible.Proxy)
	mux.Handle("GET /api/event", plausible.Proxy)

	return mux
}
