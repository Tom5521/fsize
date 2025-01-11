//go:generate go run -v ./gen-version/main.go
package meta

import (
	_ "embed"
	"strings"
)

var (
	//go:embed version.txt
	LongVersion string
	Version     = strings.SplitN(LongVersion, "-", 2)[0]
)
