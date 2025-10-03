package stat

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func processFiles(
	count, size *int64,
	path string,
	fn filepath.WalkFunc,
	mu ...*sync.Mutex,
) (err error) {
	return filepath.Walk(path,
		func(path string, info fs.FileInfo, err error) error {
			if err == nil {
				hasMutex := len(mu) > 0
				if hasMutex {
					mu[0].Lock()
				}
				*count++
				*size += info.Size()
				if hasMutex {
					mu[0].Unlock()
				}
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
