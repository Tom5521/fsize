package stat

import (
	"io/fs"
	"os"
	"path/filepath"
)

type FileLink struct {
	IsLink       bool
	PhysicalPath string
}

func NewFileLink(linfo os.FileInfo, path string) (fl FileLink, err error) {
	fl.IsLink = linfo.Mode()&fs.ModeSymlink != 0

	if fl.IsLink {
		fl.PhysicalPath, err = filepath.EvalSymlinks(path)
	}

	return
}
