package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"golidraw/rough"
)

type ExcalidrawData struct {
	Type     string              `json:"type"`
	Version  int                 `json:"version"`
	Source   string              `json:"source"`
	Elements []ExcalidrawElement `json:"elements"`
}

type ExcalidrawElement struct {
	ID              string                `json:"id"`
	Type            ExcalidrawElementType `json:"type"`
	Width           float64               `json:"width"`
	Height          float64               `json:"height"`
	X               float64               `json:"x"`
	Y               float64               `json:"y"`
	StrokeWidth     int                   `json:"strokeWidth"`
	StrokeColor     string                `json:"strokeColor"`
	BackgroundColor string                `json:"backgroundColor"`
	FrameID         string                `json:"frameId"`

	Text          string  `json:"text"`
	FontFamily    int     `json:"fontFamily"`
	FontSize      float64 `json:"fontSize"`
	TextAlign     string  `json:"textAlign"`
	VerticalAlign string  `json:"verticalAlign"`
}

type ExcalidrawElementType string

const (
	Rectangle ExcalidrawElementType = "rectangle"
	Frame     ExcalidrawElementType = "frame"
)

func Parse() (ExcalidrawData, error) {
	path := "./test.excalidraw"
	f, err := os.Open(path)
	if err != nil {
		return ExcalidrawData{}, err
	}

	var data ExcalidrawData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return ExcalidrawData{}, err
	}

	return data, nil
}

func mapToStyleAttr(style map[string]string) string {
	var sb strings.Builder
	sb.WriteString(" style=")
	sb.WriteRune('"')
	for k, v := range style {
		sb.WriteString(k)
		sb.WriteRune(':')
		sb.WriteString(v)
		sb.WriteRune(';')
	}
	sb.WriteRune('"')
	return sb.String()
}

type RenderContext struct {
	Frame   ExcalidrawElement
	Scale   float64
	OffsetX float64
	OffsetY float64
}

func renderStyle(c RenderContext, el ExcalidrawElement) map[string]string {
	res := make(map[string]string)
	res["width"] = fmt.Sprintf("%fpx", c.Scale*el.Width)
	res["height"] = fmt.Sprintf("%fpx", c.Scale*el.Height)
	res["top"] = fmt.Sprintf("%fpx", (el.Y-c.Frame.Y)*c.Scale+c.OffsetY)
	res["left"] = fmt.Sprintf("%fpx", (el.X-c.Frame.X)*c.Scale+c.OffsetX)

	if el.Type == "text" {
		res["color"] = el.StrokeColor
		if el.TextAlign == "center" {
			res["align-items"] = "center"
			res["text-align"] = "center"
		}
		if el.VerticalAlign == "middle" {
			res["justify-content"] = "center"
		}
		if el.FontFamily == 5 {
			res["font-family"] = "Excalifont"
		}
		res["font-size"] = fmt.Sprintf("%fpx", el.FontSize*c.Scale*0.9)
	}

	return res
}

func renderElement(c RenderContext, el ExcalidrawElement) string {
	switch el.Type {
	case Rectangle:
		return renderRectangle(c, el)
	case "text":
		return renderText(c, el)
	}
	panic("element type not supported")
}

func renderRectangle(c RenderContext, el ExcalidrawElement) string {
	test := rough.Svg(rough.NewRectangle(el.Width, el.Height, el.StrokeWidth))

	return fmt.Sprintf(`
<div class="element" %s>%s</div>
`, mapToStyleAttr(renderStyle(c, el)), test)
}

func renderText(c RenderContext, el ExcalidrawElement) string {
	return fmt.Sprintf(`
<div class="element" %s>%s</div>
`, mapToStyleAttr(renderStyle(c, el)), el.Text)
}

func Render(frame ExcalidrawElement, elements []ExcalidrawElement) string {
	if frame.Type != Frame {
		panic("element is not a frame")
	}

	renderHeight := 900.0
	renderWidth := 1600.0

	frameScaleRatio := renderWidth / frame.Width
	offsetY := (renderHeight - frame.Height*frameScaleRatio) / 2
	offsetX := 0.0
	if frame.Height < frame.Width {
		frameScaleRatio = renderHeight / frame.Height
		offsetY = 0
		offsetX = (renderWidth - frame.Width*frameScaleRatio) / 2
	}

	c := RenderContext{
		Frame:   frame,
		Scale:   frameScaleRatio,
		OffsetX: offsetX,
		OffsetY: offsetY,
	}

	tm := template.Must(template.New("").Funcs(template.FuncMap{
		"render": func(el ExcalidrawElement) template.HTML {
			return template.HTML(renderElement(c, el))
		},
	}).Parse(`
<!doctype html>
<html lang="en">
<head>
<style>
@font-face {
  font-family: 'Excalifont';
  font-weight: 400;
  font-style: normal;
  src: url("Excalifont-Regular.woff2") format('woff2');
}

html, body {
	padding: 0;
	margin: 0;
	background: black;
}
.content {
	background: white;
	height: 900px;
	width: 1600px;
	position: relative;
	transform-origin: top left;
}
.element {
	position: absolute;
	display: flex;
	white-space: pre;
}
</style>
</head>
<body>
<div class="content">
{{ range .Elements }}
{{ render . }}
{{ end }}
</div>
<script>
function scale() {
	const content = document.querySelector(".content");
	const s = window.innerWidth/content.clientWidth;
	content.style.transform = "scale("+s.toString()+")";
}
window.addEventListener("resize", scale);
scale();
</script>
</body>
</html>
		`))

	var sb strings.Builder
	err := tm.Execute(&sb, struct {
		Elements []ExcalidrawElement
	}{
		Elements: elements,
	})
	if err != nil {
		panic(err)
	}

	return sb.String()
}

func main() {
	data, err := Parse()
	if err != nil {
		panic(err)
	}

	var frame ExcalidrawElement
	for _, el := range data.Elements {
		if el.Type == Frame {
			frame = el
			break
		}
	}

	var elements []ExcalidrawElement
	for _, el := range data.Elements {
		if el.FrameID == frame.ID {
			elements = append(elements, el)
		}
	}

	page := Render(frame, elements)
	fmt.Println(page)
}
