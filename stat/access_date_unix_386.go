//go:build unix && !darwin && (arm || 386)
// +build unix
// +build !darwin
// +build arm 386

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
	t = time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	return
}
