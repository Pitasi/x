package togo

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

type RouteNode struct {
	Path      string
	Handler   http.HandlerFunc
	Children  map[string]*RouteNode
	IsDynamic bool
	ParamName string // only applies if IsDynamic
}

func NewRouteNode(path string, handler http.HandlerFunc) *RouteNode {
	var paramName string
	if strings.HasPrefix(path, ":") {
		paramName = path[1:]
	}

	return &RouteNode{
		Path:      path,
		Handler:   handler,
		Children:  make(map[string]*RouteNode),
		IsDynamic: paramName != "",
		ParamName: paramName,
	}
}

func (n *RouteNode) AddChild(segment string, handler http.HandlerFunc) *RouteNode {
	if c, found := n.Children[segment]; found {
		return c
	}
	c := NewRouteNode(segment, handler)
	n.Children[segment] = c
	return c
}

type RouteTree struct {
	Root *RouteNode
}

func NewRouteTree() *RouteTree {
	return &RouteTree{
		Root: NewRouteNode("/", nil),
	}
}

func (t *RouteTree) AddRoute(path string, handler http.HandlerFunc) {
	segments := strings.Split(path, "/")
	currentNode := t.Root
	for _, seg := range segments {
		if seg == "" {
			// ignore extra slashes
			continue
		}
		currentNode = currentNode.AddChild(seg, nil)
	}
	currentNode.Handler = handler
}

func (t *RouteTree) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	registerMux(mux, nil, t.Root)
	return mux
}

func registerMux(mux *http.ServeMux, segments []string, n *RouteNode) {
	if n.Path != "/" {
		if n.IsDynamic {
			segments = append(segments, "{"+n.ParamName+"}")
		} else {
			segments = append(segments, n.Path)
		}
	}

	if n.Handler != nil {
		var pattern string
		switch {
		case n.Path == "/":
			pattern = "GET /{$}"
		case n.Path == "static":
			pattern = "GET /static/"
		default:
			pattern = "GET /" + strings.Join(segments, "/")
		}
		fmt.Println("[route]", pattern)
		mux.HandleFunc(pattern, n.Handler)
	}

	for _, c := range n.Children {
		registerMux(mux, segments, c)
	}
}

func NewRouteTreeFromFS(ffs fs.FS) (*RouteTree, error) {
	tree := NewRouteTree()

	err := fs.WalkDir(ffs, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasPrefix(path, "static/") {
			return nil
		}

		if d.Type().IsRegular() &&
			!d.Type().IsDir() &&
			strings.HasSuffix(path, ".html") &&
			!strings.HasSuffix(path, "_layout.html") {
			return registerFileHTML(tree, ffs, path)
		}

		if d.Type().IsRegular() &&
			!d.Type().IsDir() &&
			strings.HasSuffix(path, ".md") {
			return registerFileMD(tree, ffs, path)
		}

		if d.Type().IsRegular() &&
			!d.Type().IsDir() &&
			strings.HasSuffix(path, ".slides") {
			return registerFileSlides(tree, ffs, path)
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	tree.AddRoute("/static/", handleStaticFiles(ffs))

	return tree, nil
}

func exists(ffs fs.FS, path string) bool {
	f, err := ffs.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}
