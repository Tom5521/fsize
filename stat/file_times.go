package stat

import (
	"os"
	"time"
)

type FileTimes struct {
	ModTime             time.Time
	AccessTime          time.Time
	CreationTime        time.Time
	SupportCreationDate bool
}

func NewFileTimes(info os.FileInfo) (times FileTimes, err error) { return newFileTimes(info) }
