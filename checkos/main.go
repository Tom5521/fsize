package checkos

import "runtime"

const (
	Linux   = runtime.GOOS == "linux"
	Windows = runtime.GOOS == "windows"
	Darwin  = runtime.GOOS == "darwin"
)
