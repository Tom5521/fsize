package meta

import (
	_ "embed"
	"strings"
)

var (
	//go:embed version.txt
	LongVersion string
	Version     string
)

func init() {
	LongVersion = strings.Replace(LongVersion, "\n", "", 1) // Cut newline
	Version = strings.SplitN(LongVersion, "-", 2)[0]
}
