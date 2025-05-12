package fsx

import (
	"io/fs"
	"os"
)

func Or(devmode bool, embedFS fs.FS, dir string) fs.FS {
	if devmode {
		root, err := os.OpenRoot(dir)
		if err != nil {
			panic(err)
		}
		return root.FS()
	}

	return embedFS
}
