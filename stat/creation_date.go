package stat

import (
	"os"
	"time"
)

func CreationDate(info os.FileInfo) (time.Time, error) { return creationDate(info) }
