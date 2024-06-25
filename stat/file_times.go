package stat

import (
	"time"
)

type FileTimes struct {
	ModTime             time.Time
	AccessTime          time.Time
	CreationTime        time.Time
	SupportCreationDate bool
}
