package stat

import (
	"os"
	"time"
)

func AccessDate(info os.FileInfo) (t time.Time, err error) { return accessDate(info) }
