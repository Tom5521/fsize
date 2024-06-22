package checkos

import "runtime"

const (
	Linux   = runtime.GOOS == "linux"
	Windows = runtime.GOOS == "windows"
	Darwin  = runtime.GOOS == "darwin"
)

var (
	UnixSystems = [...]string{
		"aix",
		"android",
		"darwin",
		"dragonfly",
		"freebsd",
		"illumos",
		"ios",
		"linux",
		"netbsd",
		"openbsd",
		"solaris",
	}
	Unix bool
)

func init() {
	for _, os := range UnixSystems {
		if os == runtime.GOOS {
			Unix = true
			break
		}
	}
}
