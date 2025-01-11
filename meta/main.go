//go:generate go run -v ./gen-version/main.go
package meta

import _ "embed"

//go:embed version.txt
var Version string
