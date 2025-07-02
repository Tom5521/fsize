//go:build unix
// +build unix

package stat

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func CreationDate(info os.FileInfo) (time.Time, error) {
	cmd := exec.Command("stat", "-c", "%W", info.Name())
	data, err := cmd.CombinedOutput()
	if err != nil {
		return time.Time{}, err
	}

	date := string(data)
	date = strings.ReplaceAll(date, "\x0a", "") // Clean stat output.

	sec, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(sec, 0), nil
}
