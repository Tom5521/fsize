//go:build unix && !darwin && (amd64 || arm64)
// +build unix
// +build !darwin
// +build amd64 arm64

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
	t = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
	return
}
