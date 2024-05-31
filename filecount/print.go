package filecount

import (
	"io/fs"

	"path/filepath"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
)

var PrintOnWalk *bool

func Print(count, size *int64, path string) (err error) {
	err = filepath.Walk(path, func(name string, info fs.FileInfo, err error) error {
		if err != nil {
			Warning(err)
			return nil
		}
		if *PrintOnWalk {
			msg.Infof("Reading \"%s\"...", name)
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
