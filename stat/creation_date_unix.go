//go:build unix
// +build unix

package stat

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

func CreationDate(info os.FileInfo) (t time.Time, err error) {
	cmd := exec.Command("stat", "-c", "%w", info.Name())
	var out []byte
	out, err = cmd.Output()
	if err != nil {
		return
	}
	date := string(out)
	if strings.HasSuffix(date, "\x0a") {
		date, _ = strings.CutSuffix(date, "\x0a")
	}
	t, err = parseStatDate(date)

	return
}

func parseStatDate(date string) (time.Time, error) {
	// EXAMPLE INPUT: "2024-06-16 21:01:08.044029927 -0400"
	const layout = "2006-01-02 15:04:05.999999999 -0700"
	return time.Parse(layout, date)
}
