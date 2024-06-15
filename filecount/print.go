package filecount

import (
	"io/fs"

	"path/filepath"

	"github.com/Tom5521/fsize/flags"
	"github.com/gookit/color"
)

func Print(count, size *int64, path string) (err error) {
	err = filepath.Walk(path, func(name string, info fs.FileInfo, err error) error {
		if err != nil {
			warning(err)
			return nil
		}
		if flags.PrintOnWalk {
			color.Infof("Reading \"%s\"...", name)
		}

		*size += info.Size()
		*count++

		return nil
	})
	if err != nil {
		return
	}
	return
}
