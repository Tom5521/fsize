package update

import "os"

func isMaybeRunningInNixOS() bool {
	_, err := os.Stat("/etc/NIXOS")
	return err == nil
}
