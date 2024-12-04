package md

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var GosmicMarkdownExtension = &gosmicMarkdownExtension{}

type gosmicMarkdownExtension struct{}

var _ (goldmark.Extender) = &gosmicMarkdownExtension{}

func (e *gosmicMarkdownExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newHTMLRenderer(), 200),
	))
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(newInlineParser(), 999),
		),
	)
}

type gosmicInlineParser struct {
	state parserState
}

type parserState int

const (
	readingTagName parserState = iota
	readingContent
)

func (g *gosmicInlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if !bytes.HasPrefix(line, []byte("::: ")) {
		return nil
	}
	block.Advance(4) // consume the "::: "

	g.state = readingTagName
	tagStart := len("::: ") // start reading after the "::: "
	cur := tagStart
	for ; cur < len(line); cur++ {
		c := line[cur]
		if !(util.IsAlphaNumeric(c)) {
			break
		}
		block.Advance(1)
	}
	tagName := line[tagStart:cur]

	for ; cur < len(line); cur++ {
		if line[cur] == ' ' {
			block.Advance(1)
		} else {
			break
		}
	}

	attrs := make(map[string]string)
	for cur < len(line) {
		attrNameStart := cur
		for ; cur < len(line); cur++ {
			c := line[cur]
			if !(util.IsAlphaNumeric(c) || c == '-' || c == '_') {
				break
			}
			block.Advance(1)
		}
		attrName := string(line[attrNameStart:cur])

		if line[cur] != '=' {
			break
		}
		cur++
		block.Advance(1)

		attrValueStart := cur
		for ; cur < len(line); cur++ {
			c := line[cur]
			if !(util.IsAlphaNumeric(c) || c == '-' || c == '_' || c == '"') {
				break
			}
			block.Advance(1)
		}
		attrValue := string(line[attrValueStart:cur])

		attrs[attrName] = attrValue

		for ; cur < len(line); cur++ {
			if line[cur] == ' ' {
				block.Advance(1)
			} else {
				break
			}
		}
	}

	block.AdvanceLine()

	var lines []string
	for {
		line, _ := block.PeekLine()
		l := strings.TrimSpace(string(line))
		if l == ":::" || len(l) == 0 {
			block.AdvanceLine()
			break
		}
		lines = append(lines, l)
		block.AdvanceLine()
	}
	content := strings.Join(lines, "\n")

	return NewDialogNode(string(tagName), attrs, string(content))
}

func (g *gosmicInlineParser) Trigger() []byte {
	return []byte(":")
}

func newInlineParser() parser.InlineParser {
	return &gosmicInlineParser{}
}

type DialogNode struct {
	ast.BaseBlock
	Tagname string
	Attrs   map[string]string
	Content string
}

// Dump implements Node.Dump.
func (n *DialogNode) Dump(source []byte, level int) {
	m := map[string]string{
		"Content": n.Content,
	}
	ast.DumpHelper(n, source, level, m, nil)
}

var KindDialog = ast.NewNodeKind("DialogNode")

func (n *DialogNode) Kind() ast.NodeKind {
	return KindDialog
}

func NewDialogNode(tagname string, attrs map[string]string, content string) *DialogNode {
	return &DialogNode{
		Tagname: tagname,
		Attrs:   attrs,
		Content: content,
	}
}

func newHTMLRenderer() renderer.NodeRenderer {
	return &dialogRenderer{}
}

type dialogRenderer struct{}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *dialogRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindDialog, r.renderDialog)
}

func (r *dialogRenderer) renderDialog(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	dialogN := node.(*DialogNode)
	switch dialogN.Tagname {
	case "dialog":
		r.writeDialog(w, dialogN.Attrs["character"], dialogN.Attrs["pos"], dialogN.Attrs["cls"], dialogN.Content)
	default:
		panic(fmt.Errorf("unknown tagname %s", dialogN.Tagname))
	}
	return ast.WalkContinue, nil
}

func (r *dialogRenderer) writeDialog(w util.BufWriter, character, pos, classes, content string) {
	w.WriteString(`<div class="article-dialog">`)
	switch pos {
	case "left":
		r.writeCharacter(w, character)
		r.writeContent(w, content, classes)
	case "right":
		r.writeContent(w, content, classes)
		r.writeCharacter(w, character)
	default:
		panic(fmt.Errorf("unknown pos %s", pos))
	}
	w.WriteString(`</div>`)
}

func (r *dialogRenderer) writeCharacter(w util.BufWriter, character string) {
	w.WriteString(`<div class="dialog-character">`)
	switch character {
	case "raisehand":
		w.WriteString(`<svg width="512" height="512" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg"><path fill="#724f3d" d="M24 7a9.83 9.83 0 0 1 2.44.31c1.86.42 4.28 1.12 6.47.7a4.2 4.2 0 0 1-.84 3.21a9.81 9.81 0 0 1 1.75 5.6v5.14H14.18v-5.17A9.82 9.82 0 0 1 24 7Z"></path><path fill="#a86c4d" d="M24 7a9.83 9.83 0 0 1 2.44.31c1.86.42 4.28 1.12 6.47.7a4.47 4.47 0 0 1-.68 3a19.4 19.4 0 0 1-5.79-.79a9.81 9.81 0 0 0-12.26 9.51v-2.94A9.82 9.82 0 0 1 24 7Z"></path><path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M24 7a9.83 9.83 0 0 1 2.44.31c1.86.42 4.28 1.12 6.47.7a4.2 4.2 0 0 1-.84 3.21a9.81 9.81 0 0 1 1.75 5.6v5.14H14.18v-5.17A9.82 9.82 0 0 1 24 7Z"></path><path fill="#ffe500" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M13.13 21.18a1.83 1.83 0 0 1 1.39-1.54l.59-.16a2.06 2.06 0 0 0 1.5-2v-2.04A1.88 1.88 0 0 1 18 13.63a15.34 15.34 0 0 0 6 1.09a15.34 15.34 0 0 0 6-1.09a1.88 1.88 0 0 1 1.4 1.81v2.06a2.06 2.06 0 0 0 1.5 2l.59.16a1.83 1.83 0 0 1 1.39 1.54a1.81 1.81 0 0 1-1.81 2H33a9 9 0 0 1-17.9 0h-.11a1.81 1.81 0 0 1-1.86-2.02Z"></path><path fill="#45413c" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M29.43 21a.77.77 0 1 1-.77-.77a.76.76 0 0 1 .77.77Zm-10.86 0a.77.77 0 1 0 .77-.77a.76.76 0 0 0-.77.77Z"></path><path fill="#ff6242" d="M26.84 25.66a.44.44 0 0 1 .33.16a.42.42 0 0 1 .1.35a3.32 3.32 0 0 1-6.54 0a.42.42 0 0 1 .1-.35a.42.42 0 0 1 .33-.16Z"></path><path fill="#ffa694" d="M24 27a4 4 0 0 1 2.52.77a3.36 3.36 0 0 1-5 0A4 4 0 0 1 24 27Z"></path><path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M26.84 25.66a.44.44 0 0 1 .33.16a.42.42 0 0 1 .1.35a3.32 3.32 0 0 1-6.54 0a.42.42 0 0 1 .1-.35a.42.42 0 0 1 .33-.16Z"></path><path fill="#ffaa54" d="M28.94 24.25a1 .6 0 1 0 2 0a1 .6 0 1 0-2 0Zm-11.88 0a1 .6 0 1 0 2 0a1 .6 0 1 0-2 0Z"></path><path fill="#ffe500" d="M16.67 10.7a1.1 1.1 0 0 0-.66-2h-.79a.28.28 0 0 1-.29-.31l.23-2.31a1.75 1.75 0 0 0-2.69-1.69a2.12 2.12 0 0 1-1.47.44a1.26 1.26 0 0 0-1.37 1a35.55 35.55 0 0 0-.73 8.66l4.86.35l.12-2.07Z"></path><path fill="#fff48c" d="M10.55 6.75a2.48 2.48 0 0 0 1.68-.43a2 2 0 0 1 2.84.68l.09-.89a1.75 1.75 0 0 0-2.69-1.69a2.12 2.12 0 0 1-1.47.41a1.26 1.26 0 0 0-1.37 1c-.08.43-.14.85-.21 1.27a1.43 1.43 0 0 1 1.13-.35Z"></path><path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M16.67 10.7a1.1 1.1 0 0 0-.66-2h-.79a.28.28 0 0 1-.29-.31l.23-2.31a1.75 1.75 0 0 0-2.69-1.69a2.12 2.12 0 0 1-1.47.44a1.26 1.26 0 0 0-1.37 1a35.55 35.55 0 0 0-.73 8.66l4.86.35l.12-2.07Z"></path><path fill="#45413c" d="M9 45.5a15 1.5 0 1 0 30 0a15 1.5 0 1 0-30 0Z" opacity=".15"></path><path fill="#00b8f0" d="M24.15 31.22a11.8 11.8 0 0 0-4.18.7a3.49 3.49 0 0 1-4.58-2.19a37 37 0 0 1-1.17-15.61H8.16a97.42 97.42 0 0 0 0 13.8a53.2 53.2 0 0 0 4 17.08h23.72v-1.51a12.14 12.14 0 0 0-11.73-12.27Z"></path><path fill="#009fd9" d="M24.15 31.22a11.8 11.8 0 0 0-4.18.7a3.49 3.49 0 0 1-4.58-2.19a33.64 33.64 0 0 1-1.49-8.79A35.75 35.75 0 0 0 15.39 33A3.49 3.49 0 0 0 20 35.19a11.8 11.8 0 0 1 4.18-.7A12.07 12.07 0 0 1 35.75 45h.13v-1.51a12.14 12.14 0 0 0-11.73-12.27Z"></path><path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M24.15 31.22a11.8 11.8 0 0 0-4.18.7a3.49 3.49 0 0 1-4.58-2.19a37 37 0 0 1-1.17-15.61H8.16a97.42 97.42 0 0 0 0 13.8a53.2 53.2 0 0 0 4 17.08h23.72v-1.51a12.14 12.14 0 0 0-11.73-12.27Z"></path><path fill="#ffe500" d="M24 30.86a8.86 8.86 0 0 1-2.54-.37v2.14a2.54 2.54 0 1 0 5.08 0v-2.14a8.86 8.86 0 0 1-2.54.37Z"></path><path fill="#ebcb00" d="M24 30.86a8.86 8.86 0 0 1-2.54-.37a2.54 2.54 0 0 0 5.08 0a8.86 8.86 0 0 1-2.54.37Z"></path><path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M24 30.86a8.86 8.86 0 0 1-2.54-.37v2.14a2.54 2.54 0 1 0 5.08 0v-2.14a8.86 8.86 0 0 1-2.54.37Z"></path></svg>`)
	case "facepalm":
		w.WriteString(`<svg width="512" height="512" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg" > <path fill="#45413c" d="M10.5 45.5a13.5 1.5 0 1 0 27 0a13.5 1.5 0 1 0-27 0Z" opacity=".15" /> <path fill="#00b8f0" d="M24 31.22A11.88 11.88 0 0 1 35.88 43.1V45H12.12v-1.9A11.88 11.88 0 0 1 24 31.22Z" /> <path fill="#009fd9" d="M24 31.22A11.88 11.88 0 0 0 12.12 43.1V45h.07A11.88 11.88 0 0 1 24 34.43A11.88 11.88 0 0 1 35.81 45h.07v-1.9A11.88 11.88 0 0 0 24 31.22Z" /> <path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M24 31.22h0A11.88 11.88 0 0 1 35.88 43.1V45h0h-23.76h0v-1.9A11.88 11.88 0 0 1 24 31.22ZM17.06 45v-2.41M30.94 45v-2.41" /> <path fill="#ebcb00" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M21.46 27.17h5.09v7.94h-5.09Z" /> <path fill="#724f3d" d="M24 6.91a9.83 9.83 0 0 0-2.44.31c-1.86.42-4.28 1.12-6.47.7a4.2 4.2 0 0 0 .84 3.21a9.81 9.81 0 0 0-1.75 5.6v5.14h19.64v-5.14A9.82 9.82 0 0 0 24 6.91Z" /> <path fill="#a86c4d" d="M24 6.91a9.83 9.83 0 0 0-2.44.31c-1.86.42-4.28 1.12-6.47.7a4.47 4.47 0 0 0 .68 3a19.4 19.4 0 0 0 5.79-.79a9.81 9.81 0 0 1 12.26 9.51v-2.91A9.82 9.82 0 0 0 24 6.91Z" /> <path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M24 6.91a9.83 9.83 0 0 0-2.44.31c-1.86.42-4.28 1.12-6.47.7a4.2 4.2 0 0 0 .84 3.21a9.81 9.81 0 0 0-1.75 5.6v5.14h19.64v-5.14A9.82 9.82 0 0 0 24 6.91Z" /> <path fill="#ffe500" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M34.87 21.12a1.83 1.83 0 0 0-1.39-1.54l-.59-.16a2.06 2.06 0 0 1-1.5-2v-2.04A1.88 1.88 0 0 0 30 13.57a15.34 15.34 0 0 1-6 1.09a15.34 15.34 0 0 1-6-1.09a1.88 1.88 0 0 0-1.4 1.81v2.06a2.06 2.06 0 0 1-1.5 2l-.59.16a1.83 1.83 0 0 0-1.39 1.54a1.81 1.81 0 0 0 1.81 2h.11a9 9 0 0 0 17.9 0h.11a1.81 1.81 0 0 0 1.82-2.02Z" /> <path fill="#ffaa54" d="M17.06 24.19a1 .6 0 1 0 2 0a1 .6 0 1 0-2 0Z" /> <path fill="#ffe500" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M30 28a5.69 5.69 0 0 0 .33-1.61c0-.79-.33-2.74.29-3.78a11.06 11.06 0 0 1 .64-1a.52.52 0 0 0-.11-.73h0a1.34 1.34 0 0 0-1.68.06h0a2.69 2.69 0 0 0-.85 1.45l-.12.61l-1.77-5.44a.85.85 0 0 0-1-.57h0a.85.85 0 0 0-.65 1.07l1.11 3.81l-1.94-4.18a.91.91 0 0 0-1.25-.4h0a.91.91 0 0 0-.42 1.14l1.63 4.28l-1.85-3.42a.8.8 0 0 0-1.1-.32h0a.81.81 0 0 0-.33 1.03l1.69 3.38l-1.29-1.89a.86.86 0 0 0-1.33-.12h0a.87.87 0 0 0-.2 1a61.55 61.55 0 0 0 3.38 6.48c.85 1 1.22 1.18 1.92 1.27c1.11.14 4.35-1.25 4.9-2.12Z" /> <path fill="#00b8f0" d="m28.86 45l-4.75-13.3a1 1 0 0 1 .52-1.24l4.85-2.3a1 1 0 0 1 1.34.49L38.13 45Z" /> <path fill="#4acfff" d="m24.73 33.43l3.87-1.84a1 1 0 0 1 1.35.5L35.72 45h2.41l-7.31-16.35a1 1 0 0 0-1.34-.49l-4.85 2.3a1 1 0 0 0-.52 1.24Z" /> <path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="m28.86 45l-4.75-13.3a1 1 0 0 1 .52-1.24l4.85-2.3a1 1 0 0 1 1.34.49L38.13 45Z" /> </svg>`)
	case "finger":
		w.WriteString(`<svg width="512" height="512" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg" > <path fill="#45413c" d="M13 45.5a11 1.5 0 1 0 22 0a11 1.5 0 1 0-22 0Z" opacity=".15" /> <path fill="#ffe500" d="M34.13 26.33V9a2.3 2.3 0 0 0-4.6 0v10.89a2.71 2.71 0 0 0-5.42 0a2.68 2.68 0 0 0-5.36 0v3.26a2.52 2.52 0 0 0-5 0v8.09a9.22 9.22 0 0 0 9.22 9.22h1.28a11.06 11.06 0 0 0 11.02-11.09a4.59 4.59 0 0 0-1.14-3.04Z" /> <path fill="#fff48c" d="M31.83 9.22a2.3 2.3 0 0 1 2.3 2.3V9a2.3 2.3 0 0 0-4.6 0v2.5a2.3 2.3 0 0 1 2.3-2.28ZM16.23 23.1a2.52 2.52 0 0 1 2.52 2.52v-2.5a2.52 2.52 0 0 0-5 0v2.5a2.52 2.52 0 0 1 2.48-2.52Zm5.2-3.42a2.68 2.68 0 0 1 2.68 2.68a2.71 2.71 0 0 1 5.42 0v-2.5a2.71 2.71 0 0 0-5.42 0a2.68 2.68 0 0 0-5.36 0v2.5a2.68 2.68 0 0 1 2.68-2.68Z" /> <path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M34.13 26.33V9a2.3 2.3 0 0 0-4.6 0v10.89a2.71 2.71 0 0 0-5.42 0v0a2.68 2.68 0 0 0-5.36 0v3.26a2.52 2.52 0 0 0-5 0v8.09a9.22 9.22 0 0 0 9.22 9.22h1.28a11.06 11.06 0 0 0 11.02-11.09a4.59 4.59 0 0 0-1.14-3.04Z" /> <path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M13.71 20.6h5.04v8.33h-5.04Z" /> <path fill="#fff48c" d="M35.27 29.37a4.6 4.6 0 0 0-4.6-4.6H20A1.23 1.23 0 0 0 18.75 26a4.38 4.38 0 0 0 .25 1.3a1.21 1.21 0 0 1 1-.56h11.1a4.55 4.55 0 0 1 4.17 2.63Z" /> <path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M35.27 29.37a4.6 4.6 0 0 0-4.6-4.6H20A1.23 1.23 0 0 0 18.75 26h0A4.22 4.22 0 0 0 23 30.22h5.24" /> <path fill="#ffe500" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M24.11 19.89v4.88m5.42-4.88v4.88" /> <path fill="none" stroke="#45413c" stroke-linecap="round" stroke-linejoin="round" d="M28.15 30.22h0a4.44 4.44 0 0 0-4.44 4.44V36" /> </svg>`)
	case "bulb":
		w.WriteString(`<img src="/static/images/bulb_icon.webp" loading="lazy">`)
	default:
		panic(fmt.Errorf("unknown character %s", character))
	}
	w.WriteString(`</div>`)
}

func (r *dialogRenderer) writeContent(w util.BufWriter, content, classes string) {
	classes = strings.Trim(classes, "\"")
	w.WriteString(`<div class="dialog-content" class="`)
	w.WriteString(classes)
	w.WriteString(`">`)
	w.WriteString(`<p>`)
	w.WriteString(content)
	w.WriteString(`</p>`)
	w.WriteString(`</div>`)
}
