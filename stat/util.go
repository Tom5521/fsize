package stat

import (
	"io/fs"
	"path/filepath"
)

func processFiles(count, size *int64, path string, fn filepath.WalkFunc) (err error) {
	return filepath.Walk(path,
		func(path string, info fs.FileInfo, err error) error {
			if err == nil {
				*count++
				*size += info.Size()
			}

			return fn(path, info, err)
		},
	)
}
