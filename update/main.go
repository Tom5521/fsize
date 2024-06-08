package update

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/Tom5521/fsize/echo"
	"github.com/Tom5521/fsize/meta"
	"github.com/minio/selfupdate"
	"github.com/schollz/progressbar/v3"
)

const UpdateURL string = "https://github.com/Tom5521/fsize/releases/latest"

const updateCompletionsInstructions string = `bash|--gen-bash-completion /usr/share/bash-completion/completions/fsize
fish|--gen-fish-completion /usr/share/fish/vendor_completions.d/fsize.fish
zsh|--gen-zsh-completion /usr/share/zsh/site-functions/_fsize`

func CheckUpdate() (tag string, latest bool, err error) {
	resp, err := http.Get(UpdateURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	latestURL := resp.Request.URL.String()
	parts := strings.Split(latestURL, "/")

	tag = parts[len(parts)-1]

	if tag == meta.Version {
		latest = true
	}

	return
}

func ApplyUpdate(tag string) (err error) {
	const baseURL string = "https://github.com/Tom5521/fsize/releases/download/%s/fsize-%s-%s%s"

	url := fmt.Sprintf(baseURL, tag, runtime.GOOS, runtime.GOARCH, func() (suffix string) {
		if runtime.GOOS == "windows" {
			suffix = ".exe"
		}
		return
	}())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading latest version...")

	_, err = io.Copy(io.MultiWriter(bar, &buf), resp.Body)
	if err != nil {
		return err
	}
	err = bar.Finish()
	if err != nil {
		return
	}
	msg.Info("Writing to binary...")
	err = selfupdate.Apply(&buf, selfupdate.Options{})
	if err != nil {
		return
	}

	msg.Info("Updating completions...")
	err = updateCompletions()
	if err == nil {
		msg.Info("Upgrade completed successfully")
	}

	return
}

func updateCompletions() (err error) {
	lines := strings.Split(updateCompletionsInstructions, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 2)
		shell, args := parts[0], parts[1]

		_, exists := exec.LookPath(shell)
		if exists != nil {
			echo.Warning(shell, "not found.")
			continue
		}
		cmd := exec.Command("fsize", strings.Fields(args)...)
		err = cmd.Run()
		if err != nil {
			return
		}
	}

	return
}