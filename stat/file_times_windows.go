//go:build windows
// +build windows

package stat

import "os"

func NewFileTimes(info os.FileInfo) (times FileTimes, err error) {
	times.ModTime = info.ModTime()
	times.SupportCreationDate = true

	times.AccessTime, err = AccessDate(info)
	if err != nil {
		return
	}
	times.CreationTime, err = CreationDate(info)
	return
}
