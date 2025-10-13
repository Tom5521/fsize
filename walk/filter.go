package walk

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Tom5521/fsize/flags"
)

func filter(path string, info fs.FileInfo, _ error) error {
	var (
		cmp   = regexp.MatchString
		match bool
		err   error
	)
	if flags.Wildcard {
		cmp = filepath.Match
	}

	if flags.Depth > 0 {
		baseDir, _ := filepath.Abs(".")
		if flags.Depth > 0 {
			absPath := filepath.Join(baseDir, path)
			parts := strings.Split(absPath, string(filepath.Separator))

			baseParts := strings.Split(baseDir, string(filepath.Separator))
			relativeDepth := len(parts) - len(baseParts)

			if uint(relativeDepth) >= flags.Depth && info.IsDir() {
				return filepath.SkipDir
			}
		}
	}

	for _, pattern := range flags.Patterns {
		match, err = cmp(pattern, path)
		if err != nil || !match {
			return filepath.SkipDir
		}
	}
	for _, pattern := range flags.IgnorePatterns {
		match, err = cmp(pattern, path)
		if err != nil || match {
			return filepath.SkipDir
		}
	}

	return nil
}
