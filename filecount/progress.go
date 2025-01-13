package filecount

import (
	"io/fs"
	"path/filepath"

	_ "unsafe"

	"github.com/Tom5521/fsize/echo"
	po "github.com/leonelquinteros/gotext"
	"github.com/schollz/progressbar/v3"
)

func Progress(count, size *int64, path string) (err error) {
	echo.Info(po.Get("Counting the amount of files..."))

	var warnings []error

	countBar := progressbar.Default(-1)
	setbar := func(files, errors int64) {
		countBar.Describe(
			po.Get("%v files with %v errors", files, errors),
		)
	}
	err = filepath.Walk(path, func(_ string, info fs.FileInfo, err error) error {
		if err != nil {
			warnings = append(warnings, err)
			return nil
		}
		*count++
		*size += info.Size()
		setbar(*count, int64(len(warnings)))
		return nil
	})
	if err != nil {
		warnings = append(warnings, err)
	}

	err = countBar.Finish()
	if err != nil {
		warnings = append(warnings, err)
	}

	for _, e := range warnings {
		echo.Warning(e.Error())
	}
	if err != nil {
		return
	}

	return
}
