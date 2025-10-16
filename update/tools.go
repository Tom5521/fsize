package update

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func isMaybeRunningInNixOS() bool {
	_, err := os.Stat("/etc/NIXOS")
	return err == nil
}

func binaryDigest() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}

	bin, err := os.ReadFile(executable)
	if err != nil {
		return "", err
	}

	sha := sha256.Sum256(bin)
	return "sha256:" + hex.EncodeToString(sha[:]), nil
}
