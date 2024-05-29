//go:build windows
// +build windows

package filestat

import (
	"os"
	"syscall"
	"time"
)

func CreationDate(info os.FileInfo) time.Time {
	d := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, d.CreationTime.Nanoseconds())
}
