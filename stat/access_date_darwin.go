//go:build darwin
// +build darwin

package stat

import (
	"os"
	"syscall"
	"time"
)

func AccessDate(info os.FileInfo) (t time.Time, err error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return t, ErrGettingStruct
	}
	t = time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
	return
}
