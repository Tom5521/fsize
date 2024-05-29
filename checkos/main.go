package checkos

import "runtime"

var (
	UnixSystems = []string{
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
	Unix = func() bool {
		for _, os := range UnixSystems {
			if os == runtime.GOOS {
				return true
			}
		}

		return false
	}()
)

const (
	Linux   = runtime.GOOS == "linux"
	Windows = runtime.GOOS == "windows"
	Darwin  = runtime.GOOS == "darwin"
)
