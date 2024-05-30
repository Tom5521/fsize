//go:build unix
// +build unix

package stat

import (
	"os"
	"time"
)

func CreationDate(info os.FileInfo) time.Time {
	return info.ModTime()
}
