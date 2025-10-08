//go:build unix

package stat

import (
	"os"
)

func NewFileTimes(info os.FileInfo) (times FileTimes, err error) {
	times.ModTime = info.ModTime()

	times.CreationTime, err = CreationDate(info)
	times.SupportCreationDate = err == nil

	times.AccessTime, err = AccessDate(info)
	if err != nil {
		return times, err
	}

	return times, err
}
