//go:build darwin || freebsd || netbsd

package stat

import (
	"os"
	"syscall"
	"time"
)

func accessDate(info os.FileInfo) (t time.Time, err error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return t, ErrGettingStruct
	}
	t = time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec))
	return t, err
}
