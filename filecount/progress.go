package filecount

import (
	"fmt"
	"io/fs"
	"path/filepath"

	_ "unsafe"

	_ "github.com/Tom5521/fsize/echo"
	"github.com/gookit/color"
	"github.com/schollz/progressbar/v3"
)

//go:linkname warning github.com/Tom5521/fsize/echo.Warning
func warning(...any) // FUCK THE IMPORT CYCLE.

func Progress(count, size *int64, path string) (err error) {
	color.Infoln("Counting the amount of files...")
	var warnings []error
	countBar := progressbar.Default(-1)
	setbar := func(files, errors int64) {
		countBar.Describe(
			fmt.Sprintf("%v files with %v errors", files, errors),
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
		warning(e.Error())
	}
	if err != nil {
		return
	}

	return
}
