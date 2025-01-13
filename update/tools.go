package update

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func isMaybeRunningInNixOS() bool {
	_, err := os.Stat("/etc/NIXOS")
	return err == nil
}

func isMaybeRunningInTermux() (ok bool) {
	// There are more efficient and safer ways to check this, but if either of these two do not work it means that the USER is the one who does not want the program to know that it is running on termux.
	if runtime.GOOS != "android" {
		return
	}
	prefix, ok := os.LookupEnv("PREFIX")
	if !ok {
		return
	}
	ok = strings.Contains(prefix, "com.termux")
	if !ok {
		_, err := exec.LookPath("termux-setup-storage")
		ok = err == nil
	}

	return
}
