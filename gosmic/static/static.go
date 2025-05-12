package static

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"net/http"
)

func generateStaticFileHashes(staticFS fs.FS) map[string]string {
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

type StaticFS struct {
	resources fs.FS
	hashes    map[string]string
}

func NewStaticFS(resources fs.FS) *StaticFS {
	return &StaticFS{
		resources: resources,
		hashes:    generateStaticFileHashes(resources),
	}
}

func (s *StaticFS) FileHash(name string) string {
	return s.hashes[name]
}

func (s *StaticFS) Handler() http.Handler {
	fs := http.FS(s.resources)
	h := http.FileServer(neuteredFileSystem{fs})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		h.ServeHTTP(w, r)
	})
}
