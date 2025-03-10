package togo

import (
	"io/fs"
	"net/http"
	"strings"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(name string) (http.File, error) {
	if !fs.ValidPath(strings.TrimPrefix(name, "/")) {
		return nil, fs.ErrNotExist
	}

	f, err := nfs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return nil, fs.ErrNotExist
	}

	return f, nil
}

func handleStaticFiles(staticFS fs.FS) http.HandlerFunc {
	var fs http.FileSystem = http.FS(staticFS)
	h := http.FileServer(neuteredFileSystem{fs})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		h.ServeHTTP(w, r)
	}
}
