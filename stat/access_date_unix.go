//go:build unix && !darwin
// +build unix,!darwin

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
	// Make the conversion if having the 386 architecture.
	t = time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	return
}
