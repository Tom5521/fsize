//go:build windows

package stat

import (
	"os"
	"syscall"
	"time"
)

func creationDate(info os.FileInfo) (time.Time, error) {
	d := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, d.CreationTime.Nanoseconds()), nil
}
