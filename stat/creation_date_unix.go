//go:build unix
// +build unix

package stat

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func CreationDate(info os.FileInfo) (t time.Time, err error) {
	cmd := exec.Command("stat", "-c", "%w", info.Name())
	var builder strings.Builder
	cmd.Stdout = &builder
	err = cmd.Run()
	if err != nil {
		return t, err
	}

	date := builder.String()
	date = regexp.MustCompile(`\x0a`).ReplaceAllString(date, "")
	t, err = parseStatDate(date)

	return
}

func parseStatDate(date string) (time.Time, error) {
	// EXAMPLE INPUT: "2024-06-16 21:01:08.044029927 -0400"
	const layout = "2006-01-02 15:04:05.999999999 -0700"
	return time.Parse(layout, date)
}
