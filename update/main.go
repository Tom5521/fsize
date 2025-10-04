package update

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Tom5521/fsize/checkos"
	"github.com/Tom5521/fsize/echo"
	"github.com/Tom5521/fsize/meta"
	"github.com/gookit/color"
	po "github.com/leonelquinteros/gotext"
	"github.com/schollz/progressbar/v3"
)

const UpdateURL string = "https://github.com/Tom5521/fsize/releases/latest"

func CheckUpdate() (tag string, latest bool, err error) {
	echo.Info(po.Get("Checking the latest version available..."))
	resp, err := http.Get(UpdateURL)
	if err != nil {
		err = fmt.Errorf("error getting http response: %v", err)
		return tag, latest, err
	}
	defer resp.Body.Close()

	latestURL := resp.Request.URL.String()
	parts := strings.Split(latestURL, "/")

	tag = parts[len(parts)-1]

	if tag == meta.Version {
		latest = true
	}

	return tag, latest, err
}

func ApplyUpdate(tag string) (err error) {
	const baseURL string = "https://github.com/Tom5521/fsize/releases/download/%s/fsize-%s-%s"

	// This is for the termux users.
	goos := runtime.GOOS
	if runtime.GOOS == "android" {
		goos = "linux"
	}

	url := fmt.Sprintf(
		baseURL,
		tag,
		goos,
		runtime.GOARCH,
	)
	if checkos.Windows {
		url += ".exe"
	}

	var needConfirm bool
	if isMaybeRunningInNixOS() {
		needConfirm = true
		echo.Warning(po.Get("NB! It seems that you are in a NixOS."))
		echo.Warning(po.Get(
			"Due to the non-standard filesystem implementation of the environment, the update command may not work as expected.",
		))
	}

	if needConfirm {
		var confirm bool
		prompt := &survey.Confirm{
			Message: po.Get("Do you want to proceed with the update?"),
		}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			echo.Info(po.Get("The command has been cancelled."))
			return nil
		}
	}

	executable, err := os.Executable()
	if err != nil {
		return err
	}
	oldExec := executable + ".old"

	if err = os.Rename(executable, oldExec); err != nil {
		return errors.New(po.Get("error renaming the old executable [%s] to: %v", executable, err))
	}

	resp, err := http.Get(url)
	if err != nil {
		return errors.New(po.Get("error getting http response: %v", err))
	}

	defer resp.Body.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go catchSIGINT(resp, sigchan)

	bar := progressbar.DefaultBytes(resp.ContentLength, po.Get("Downloading latest version..."))

	if err = download(executable, resp, bar); err != nil {
		echo.Error(po.Get("Error downloading binary to %s: %v", executable, err))
		echo.Warning(po.Get("Reversing changes..."))
		err = os.Rename(oldExec, executable)
		if err != nil {
			return errors.New(
				po.Get(
					"error reversing changes: failed to rename %s to %s: %v",
					oldExec,
					executable,
					err,
				),
			)
		}
		return err
	}

	signal.Stop(sigchan)
	close(sigchan)

	if err = bar.Finish(); err != nil {
		return errors.New(po.Get("error finishing the progress bar"))
	}

	if err = os.Remove(oldExec); err != nil {
		return errors.New(po.Get("error removing old executable: %v", err))
	}

	if checkos.Unix {
		echo.Info(po.Get("Updating completions..."))
		err = updateCompletions(executable)
		if err != nil {
			return fmt.Errorf("error updating the completions: %v", err)
		}
	}

	echo.Info(po.Get("Upgrade completed successfully"))
	fmt.Printf("%s -> %s\n", color.Red.Render(meta.Version), color.Green.Render(tag))

	return err
}

func catchSIGINT(
	resp *http.Response,
	sigchan chan os.Signal,
) {
	c := <-sigchan
	if c == nil {
		return
	}
	echo.Info(po.Get("%s detected, reversing changes before finishing program...", c.String()))
	resp.Body.Close()
}

func download(executable string, resp *http.Response, bar *progressbar.ProgressBar) (err error) {
	file, err := os.OpenFile(executable, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(io.MultiWriter(bar, file), resp.Body)
	if err != nil {
		return err
	}

	return err
}
