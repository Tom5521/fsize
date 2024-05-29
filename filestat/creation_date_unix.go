//go:build unix
// +build unix

package filestat

import (
	"os"
	"time"
)

func CreationDate(info os.FileInfo) time.Time {
	return info.ModTime()
}
