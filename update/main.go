package update

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/Tom5521/fsize/checkos"
	"github.com/Tom5521/fsize/meta"
	"github.com/gookit/color"
	"github.com/minio/selfupdate"
	"github.com/schollz/progressbar/v3"
)

const UpdateURL string = "https://github.com/Tom5521/fsize/releases/latest"

func CheckUpdate() (tag string, latest bool, err error) {
	color.Infoln("Checking the latest version available...")
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
	const baseURL string = "https://github.com/Tom5521/fsize/releases/download/%s/fsize-%s-%s"

	url := fmt.Sprintf(
		baseURL,
		tag,
		runtime.GOOS,
		runtime.GOARCH,
	)
	if checkos.Windows {
		url += ".exe"
	}

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
	color.Infoln("Writing to binary...")
	err = selfupdate.Apply(&buf, selfupdate.Options{})
	if err != nil {
		return
	}

	if checkos.Unix {
		color.Infoln("Updating completions...")
		err = updateCompletions()
		if err != nil {
			return
		}
	}

	color.Infoln("Upgrade completed successfully")
	fmt.Printf("%s -> %s\n", color.Red.Render(meta.Version), color.Green.Render(tag))

	return
}
