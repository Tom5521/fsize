package stat

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/Tom5521/fsize/flags"
)

func processFiles(
	count, size *int64,
	path string,
	fn filepath.WalkFunc,
	mu *sync.Mutex,
) (err error) {
	return filepath.Walk(path,
		func(path string, info fs.FileInfo, err2 error) error {
			var (
				wd    string
				rel   string
				cmp   = regexp.MatchString
				match bool
				err   error
			)
			if flags.Wildcard {
				cmp = filepath.Match
			}
			wd, err = os.Getwd()
			if err != nil {
				return err
			}
			rel, err = filepath.Rel(wd, path)
			if err != nil {
				return err
			}
			preferredPath := rel
			if strings.HasPrefix(rel, "..") {
				preferredPath = path
			}

			for _, pattern := range flags.Patterns {
				match, err = cmp(pattern, preferredPath)
				if err != nil || !match {
					return filepath.SkipDir
				}
			}
			for _, pattern := range flags.IgnorePatterns {
				match, err = cmp(pattern, preferredPath)
				if err != nil || match {
					return filepath.SkipDir
				}
			}

			if err2 == nil {
				mu.Lock()
				*count++
				*size += info.Size()
				mu.Unlock()
			}

			return fn(preferredPath, info, err)
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
