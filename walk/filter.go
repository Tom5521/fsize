package walk

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Tom5521/fsize/flags"
	"github.com/karrick/godirwalk"
)

func filter(path string, info *godirwalk.Dirent) error {
	var (
		cmp   = regexp.MatchString
		match bool
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
				return godirwalk.SkipThis
			}
		}
	}

	var err error
	for _, pattern := range flags.Patterns {
		match, err = cmp(pattern, path)
		if err != nil || !match {
			return godirwalk.SkipThis
		}
	}
	for _, pattern := range flags.IgnorePatterns {
		match, err = cmp(pattern, path)
		if err != nil || match {
			return godirwalk.SkipThis
		}
	}

	return nil
}
