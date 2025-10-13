//go:build linux

package stat

import (
	"errors"
	"os"
	"time"

	po "github.com/leonelquinteros/gotext"
	"golang.org/x/sys/unix"
)

func creationDate(info os.FileInfo) (time.Time, error) {
	var statx unix.Statx_t
	err := unix.Statx(
		unix.AT_FDCWD,
		info.Name(),
		unix.AT_SYMLINK_NOFOLLOW,
		unix.STATX_BTIME,
		&statx,
	)
	if err != nil {
		if err == unix.ENOSYS {
			err = errors.New(
				po.Get(
					"The statx syscall is not supported. At least Linux kernel 4.11 is needed",
				),
			)
		}
		return time.Time{}, err
	}
	return time.Unix(int64(statx.Btime.Sec), int64(statx.Btime.Nsec)), nil
}
