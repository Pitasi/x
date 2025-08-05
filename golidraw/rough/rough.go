package rough

import "fmt"

type Rectangle struct {
	BaseDrawable
}

func NewRectangle(width, height float64, strokeWidth int) Rectangle {
	return Rectangle{
		BaseDrawable: BaseDrawable{
			width:       width,
			height:      height,
			strokeWidth: strokeWidth,
		},
	}
}

func (r Rectangle) Ops() []Op { return nil }

type BaseDrawable struct {
	width       float64
	height      float64
	stroke      string
	strokeWidth int
	fill        string
}

func (d BaseDrawable) Width() float64   { return d.width }
func (d BaseDrawable) Height() float64  { return d.height }
func (d BaseDrawable) Stroke() string   { return d.stroke }
func (d BaseDrawable) StrokeWidth() int { return d.strokeWidth }
func (d BaseDrawable) Fill() string     { return d.fill }

type Drawable interface {
	Width() float64
	Height() float64
	Stroke() string
	StrokeWidth() int
	Fill() string
	Ops() []Op
}

type Op interface {
	Render() string
}

func Svg(d Drawable) string {
	return fmt.Sprintf(`
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 252.58984375 88.41015625" preserveAspectRatio="none" style="width:%fpx;height:%fpx;overflow:visible;">
<g stroke-linecap="round" transform="translate(10 10)">
  <path d="M17.1 0 C88.53 -0.51, 159.36 -1.98, 215.49 0 M17.1 0 C88.28 -1, 159.22 -1.95, 215.49 0 M215.49 0 C227.27 -0.97, 232.13 5.79, 232.59 17.1 M215.49 0 C225.05 -0.92, 234.55 4.38, 232.59 17.1 M232.59 17.1 C230.66 31.46, 233.5 42.33, 232.59 51.31 M232.59 17.1 C232.25 26.68, 233.18 37.74, 232.59 51.31 M232.59 51.31 C233.54 61.94, 227.18 66.69, 215.49 68.41 M232.59 51.31 C232.81 61.97, 225.76 69.81, 215.49 68.41 M215.49 68.41 C140.05 69.29, 63.71 67.73, 17.1 68.41 M215.49 68.41 C155.09 69.64, 94.3 69.12, 17.1 68.41 M17.1 68.41 C3.74 67.79, 1.39 60.84, 0 51.31 M17.1 68.41 C5.92 68.56, 0.98 62.78, 0 51.31 M0 51.31 C1.45 41.67, -0.58 31.08, 0 17.1 M0 51.31 C0.95 37.64, -0.44 24.92, 0 17.1 M0 17.1 C-1.67 7.44, 4.39 -0.86, 17.1 0 M0 17.1 C1.81 3.59, 5.46 1.17, 17.1 0" stroke="#1e1e1e" stroke-width="%d" fill="none">
  </path>
</g>
</svg>`, d.Width(), d.Height(), d.StrokeWidth())
}
