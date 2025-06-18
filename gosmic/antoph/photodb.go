package antoph

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/webp"
)

type Tag[T any] struct {
	Value T `json:"value"`
}

func (t *Tag[T]) Get() (v T) {
	if t == nil {
		return v
	}
	return t.Value
}

type photoInfo struct {
	Make             *Tag[string]
	Model            *Tag[string]
	ISO              *Tag[int]
	ExposureTime     *Tag[string]
	FNumber          *Tag[string]
	LensModel        *Tag[string]
	DateTimeOriginal *Tag[string]
	DateTime         *Tag[string]
	Keywords         *Tag[any]
}

func (info photoInfo) GetDateTime() time.Time {
	if info.DateTimeOriginal != nil {
		return parseExifDate(info.DateTimeOriginal.Value)
	}
	if info.DateTime != nil {
		return parseExifDate(info.DateTime.Value)
	}
	//panic("image doesn't have required EXIF DateTime or DateTimeOriginal tag")
	return time.Now().Add(365 * 24 * time.Hour) // TODO: fix the image instead
}

func (info photoInfo) GetKeywords() []string {
	if info.Keywords == nil {
		return nil
	}

	switch v := info.Keywords.Value.(type) {
	case string:
		return []string{v}
	case []any:
		strs := make([]string, len(v))
		for i, s := range v {
			strs[i] = s.(string)
		}
		return strs
	}

	panic(fmt.Sprintf("invalid Keywords value type: %T: %v", info.Keywords.Value, info.Keywords))
}

func (info photoInfo) GetCamera() string {
	if info.Make == nil && info.Model == nil {
		return ""
	}
	return fmt.Sprintf("%s %s", info.Make.Value, info.Model.Value)
}

func parseExifDate(s string) time.Time {
	t, err := time.Parse("2006:01:02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func extractMeta(p string) ImgMeta {
	js, err := os.Open(path.Join(p, "info.json"))
	if err != nil {
		panic(err)
	}
	defer js.Close()

	var img *os.File
	for _, name := range []string{"w_2500.webp", "w_1900.webp", "w_1200.webp"} {
		var err error
		img, err = os.Open(path.Join(p, name))
		if err != nil {
			continue
		}
		defer img.Close()
		break
	}
	if img == nil {
		panic(fmt.Sprintf("cannot find image file in %s\n", p))
	}

	imgCfg, err := webp.DecodeConfig(img)
	if err != nil {
		panic(err)
	}

	var info photoInfo
	err = json.NewDecoder(js).Decode(&info)
	if err != nil {
		panic(err)
	}

	camera := info.GetCamera()
	lens := info.LensModel.Get()
	iso := info.ISO.Get()
	aperture := parseFNumber(info.FNumber.Get())
	shutterSpeed := info.ExposureTime.Get()
	datetime := info.GetDateTime()
	keywords := info.GetKeywords()

	return ImgMeta{
		Width:        imgCfg.Width,
		Height:       imgCfg.Height,
		Date:         datetime,
		Camera:       camera,
		Lens:         lens,
		ISO:          iso,
		ShutterSpeed: shutterSpeed,
		Aperture:     aperture,
		Keywords:     keywords,
	}
}

func parseFNumber(s string) string {
	if len(s) == 0 {
		return ""
	}

	parts := strings.SplitN(s, "/", 2)
	if len(parts) == 1 {
		return s
	}

	n, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}

	d, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%.1f", float64(n)/float64(d))
}
