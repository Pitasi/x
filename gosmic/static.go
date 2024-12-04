package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

func generateStaticFileHashes(staticFS embed.FS) map[string]string {
	hashes := make(map[string]string)
	err := fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		file, err := staticFS.Open(path)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		hash, err := generateFileHash(file)
		if err != nil {
			return fmt.Errorf("generating hash for %s: %w", path, err)
		}
		hashes[path] = hash
		return nil
	})
	if err != nil {
		panic(err)
	}
	return hashes
}

func generateFileHash(r io.Reader) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)[0:3]), nil
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(name string) (http.File, error) {
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

//go:embed static
var staticFS embed.FS

var staticFileHashes = generateStaticFileHashes(staticFS)

func static(mux *http.ServeMux) {
	var fs http.FileSystem = http.FS(staticFS)
	if devMode {
		fs = http.FS(os.DirFS("./"))
	}
	h := http.FileServer(neuteredFileSystem{fs})

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		h.ServeHTTP(w, r)
	})
}
