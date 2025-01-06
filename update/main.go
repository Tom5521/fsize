package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Tom5521/fsize/checkos"
	"github.com/Tom5521/fsize/echo"
	"github.com/Tom5521/fsize/locales"
	"github.com/Tom5521/fsize/meta"
	"github.com/gookit/color"
	"github.com/schollz/progressbar/v3"
)

const UpdateURL string = "https://github.com/Tom5521/fsize/releases/latest"

var po = locales.Po

func CheckUpdate() (tag string, latest bool, err error) {
	echo.Info("Checking the latest version available...")
	resp, err := http.Get(UpdateURL)
	if err != nil {
		err = fmt.Errorf("error getting http response: %v", err)
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

	var needConfirm bool
	if isMaybeRunningInNixOS() {
		needConfirm = true
		echo.Warning("NB! It seems that you are in a NixOS.")
		echo.Warning(
			"Due to the non-standard filesystem implementation of the environment, the update command may not work as expected.",
		)
	}

	if needConfirm {
		var confirm bool
		prompt := &survey.Confirm{
			Message: po.Get("Do you want to proceed with the update?"),
		}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			echo.Info("The command has been cancelled.")
			return nil
		}
	}

	executable, err := os.Executable()
	if err != nil {
		return
	}
	oldExec := executable + ".old"

	if err = os.Rename(executable, oldExec); err != nil {
		return fmt.Errorf("error renaming the old executable [%s] to: %v", executable, err)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating a new http request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting http response: %v", err)
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(resp.ContentLength, po.Get("Downloading latest version..."))

	if err = download(executable, resp, bar); err != nil {
		echo.Error("Error downloading binary to %s: %v", executable, err)
		echo.Warning("Reversing changes...")
		err = os.Rename(oldExec, executable)
		if err != nil {
			return fmt.Errorf(
				"error reversing changes: falied to rename %s to %s: %v",
				oldExec,
				executable,
				err,
			)
		}
		return
	}

	if err = bar.Finish(); err != nil {
		return fmt.Errorf("error finishing the progress bar")
	}

	if err = os.Remove(oldExec); err != nil {
		return fmt.Errorf("error removing old executable: %v", err)
	}

	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
	if usr.Username != "root" && checkos.Unix {
		echo.Warning("The user is not root, the completions will not be updated")
	}

	if checkos.Unix && usr.Username == "root" {
		echo.Info("Updating completions...")
		err = updateCompletions(executable)
		if err != nil {
			return fmt.Errorf("error updating the completions: %v", err)
		}
	}

	echo.Info("Upgrade completed successfully")
	fmt.Printf("%s -> %s\n", color.Red.Render(meta.Version), color.Green.Render(tag))

	return
}

func download(executable string, resp *http.Response, bar *progressbar.ProgressBar) (err error) {
	file, err := os.OpenFile(executable, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(io.MultiWriter(bar, file), resp.Body)
	if err != nil {
		return err
	}

	return
}
