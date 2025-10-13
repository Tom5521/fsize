//go:build openbsd || dragonfly || solaris

package stat

import (
	"errors"
	"os"
	"time"

	po "github.com/leonelquinteros/gotext"
)

func creationDate(_ os.FileInfo) (time.Time, error) {
	return time.Time{}, errors.New(
		po.Get(
			"file birthtime is not supported",
		),
	)
}
