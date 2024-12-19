package socialimg

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Generator struct {
	font          func(points float64) font.Face
	propicResized image.Image
	base          image.Image
}

func NewGenerator(ttf, propicFile io.Reader) (*Generator, error) {
	font, err := loadFont(ttf)
	if err != nil {
		return nil, fmt.Errorf("loading font: %s", err)
	}

	propic, _, err := image.Decode(propicFile)
	if err != nil {
		return nil, fmt.Errorf("opening image file: %s", err)
	}
	width := 128
	propicResized := imaging.Resize(propic, width, 0, imaging.Box)

	gen := &Generator{
		font:          font,
		propicResized: propicResized,
	}
	gen.generateBase()

	return gen, nil
}

func (gen *Generator) generateBase() {
	dc := gg.NewContext(1024, 587)

	bgColor := "#1e1a4d"
	dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.SetHexColor(bgColor)
	dc.Fill()

	g := gg.NewLinearGradient(0, float64(dc.Height()), float64(dc.Width()), 0)
	g.AddColorStop(0, color.RGBA{
		R: 0xF5,
		G: 0xB3,
		B: 0xFF,
		A: 0xFF,
	})
	g.AddColorStop(0.5, color.RGBA{
		R: 0xF5,
		G: 0xB3,
		B: 0xFF,
		A: 0xFF,
	})
	g.AddColorStop(1, color.RGBA{
		R: 0xFB,
		G: 0xDE,
		B: 0xFF,
		A: 0xFF,
	})

	// background
	border := 10.0
	dc.DrawRoundedRectangle(border, border, float64(dc.Width())-(2*border), float64(dc.Height())-(2*border), 20)
	dc.SetFillStyle(g)
	dc.Fill()

	dc.DrawImage(gen.propicResized, 45, 40)

	gen.base = dc.Image()
}

func (gen *Generator) Generate(w io.Writer, title, subtitle string) error {
	dc := gg.NewContext(1024, 587)
	dc.DrawImage(gen.base, 0, 0)

	textColor := "#1e1a4d"
	fontSize := 65.0
	dc.SetFontFace(gen.font(fontSize))
	textRightMargin := 55.0
	textTopMargin := 290.0
	x := textRightMargin
	y := textTopMargin
	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin
	dc.SetHexColor(textColor)
	dc.DrawStringWrapped(title, x, y, 0, 0, maxWidth, 1, gg.AlignLeft)
	txtLines := dc.WordWrap(title, maxWidth)
	_, titleH := dc.MeasureString(title)
	titleH = titleH * float64(len(txtLines))

	dateColor := "#00000070"
	dateFontSize := 35.0
	dc.SetFontFace(gen.font(dateFontSize))
	dc.SetHexColor(dateColor)
	dateMarginTop := 50.
	dateY := y + titleH + dateMarginTop
	dc.DrawString(subtitle, textRightMargin, dateY)

	return jpeg.Encode(w, dc.Image(), nil)
}

func loadFont(r io.Reader) (func(points float64) font.Face, error) {
	fontBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	return func(points float64) font.Face {
		return truetype.NewFace(f, &truetype.Options{
			Size: points,
		})
	}, nil
}
