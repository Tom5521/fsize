package filecount

import (
	"fmt"
	"io/fs"
	"path/filepath"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/schollz/progressbar/v3"
)

var Warning func(...any)

func Progress(count, size *int64, path string) (err error) {
	msg.Info("Counting the amount of files...")
	var warnings []error
	countBar := progressbar.Default(-1)
	setbar := func(files, errors int) {
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
		setbar(int(*count), len(warnings))
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
		Warning(e.Error())
	}
	if err != nil {
		return
	}

	return
}
