//go:build windows
// +build windows

package stat

import (
	"os"
	"syscall"
	"time"
)

func AccessDate(info os.FileInfo) (t time.Time, err error) {
	stat, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return t, ErrGettingStruct
	}

	t = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	return
}
