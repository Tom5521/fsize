//go:build windows
// +build windows

package stat

import (
	"os"
	"syscall"
	"time"
)

func CreationDate(info os.FileInfo) (time.Time, error) {
	d := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, d.CreationTime.Nanoseconds()), nil
}
