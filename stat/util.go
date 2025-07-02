package stat

import (
	"io/fs"
	"os"
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

type absFileInfo struct {
	absPath string
	os.FileInfo
}

func newAbsFileInfo(info os.FileInfo, absPath string) absFileInfo {
	return absFileInfo{absPath: absPath, FileInfo: info}
}

// Returns the absolute path of the file.
func (cf absFileInfo) Name() string {
	return cf.absPath
}
