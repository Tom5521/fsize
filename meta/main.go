package meta

import (
	"strings"
)

var (
	LongVersion string = "go_installed"
	Version     string
)

func init() {
	Version = strings.SplitN(LongVersion, "-", 2)[0]
}
