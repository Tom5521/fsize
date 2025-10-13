//go:build freebsd || netbsd || darwin

package stat

import (
	"errors"
	"os"
	"time"

	po "github.com/leonelquinteros/gotext"
	"golang.org/x/sys/unix"
)

func creationDate(info os.FileInfo) (time.Time, error) {
	var stat unix.Stat_t
	err := unix.Stat(
		info.Name(),
		&stat,
	)
	if err != nil {
		if err == unix.ENOSYS {
			err = errors.New(
				po.Get(
					"file birthtime is not supported",
				),
			)
		}
		return time.Time{}, err
	}
	return time.Unix(int64(stat.Btim.Sec), int64(stat.Btim.Nsec)), nil
}
