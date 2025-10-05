package stat

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/Tom5521/fsize/flags"
)

func processFiles(
	count, size *int64,
	path string,
	fn filepath.WalkFunc,
	mu ...*sync.Mutex,
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

			if flags.Pattern != "" {
				match, err = cmp(flags.Pattern, rel)
				if err != nil || !match {
					return err
				}
			}
			if flags.IgnorePattern != "" {
				match, err = cmp(flags.IgnorePattern, rel)
				if err != nil || match {
					return err
				}
			}
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
