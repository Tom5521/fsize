package update

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/Tom5521/fsize/checkos"
	"github.com/Tom5521/fsize/meta"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v76/github"
	"github.com/gookit/color"
	po "github.com/leonelquinteros/gotext"
	"github.com/schollz/progressbar/v3"
)

const UpdateURL string = "https://github.com/Tom5521/fsize/releases/latest"

type Status struct {
	Asset         *github.ReleaseAsset
	Release       *github.RepositoryRelease
	Latest        bool
	CurrentDigest string
}

func CheckUpdate(client *github.Client) (status *Status, err error) {
	log.Info(po.Get("Checking the latest version available..."))
	rel, _, err := client.Repositories.GetLatestRelease(context.Background(), "Tom5521", "fsize")
	if err != nil {
		return nil, err
	}

	binaryName := fmt.Sprintf("fsize-%s-%s",
		runtime.GOOS,
		runtime.GOARCH,
	)
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	var relAsset *github.ReleaseAsset
	for _, ra := range rel.Assets {
		if ra.GetName() == binaryName {
			relAsset = ra
			break
		}
	}
	if relAsset == nil {
		return &Status{
			Release: rel,
			Latest:  false,
		}, nil
	}

	binDigest, err := binaryDigest()
	if err != nil {
		return nil, err
	}

	return &Status{
		CurrentDigest: binDigest,
		Asset:         relAsset,
		Release:       rel,
		Latest: rel.GetTagName() == meta.Version &&
			binDigest == relAsset.GetDigest(),
	}, nil
}

func ApplyUpdate(client *github.Client, status *Status) (err error) {
	if status.Asset == nil {
		return errors.New(
			po.Get(`unsupported platform, try building fsize by yourself`),
		)
	}
	var needConfirm bool
	if isMaybeRunningInNixOS() {
		needConfirm = true
		log.Warn(po.Get("NB! It seems that you are in a NixOS."))
		log.Warn(po.Get(
			"Due to the non-standard filesystem implementation of the environment, the update command may not work as expected.",
		))
	}

	if needConfirm {
		var confirm bool
		err = huh.NewConfirm().
			Title(po.Get("Do you want to proceed with the update?")).
			Affirmative(po.Get("Yes")).
			Negative("No").
			Value(&confirm).
			WithTheme(huh.ThemeBase()).
			WithAccessible(true).
			Run()

		if !confirm || err != nil {
			if err != nil {
				log.Error(po.Get("error running confirm dialog"), "err", err)
			} else {
				log.Info(po.Get("The command has been cancelled."))
			}
			return err
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

	resp, err := http.Get(status.Asset.GetBrowserDownloadURL())
	if err != nil {
		return errors.New(po.Get("error getting http response: %v", err))
	}

	defer resp.Body.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go catchSIGINT(resp, sigchan)

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetDescription(po.Get("Downloading latest version...")),
		progressbar.OptionShowBytes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetVisibility(true),
	)

	if err = download(executable, resp, bar); err != nil {
		log.Error(po.Get("Error downloading binary to %s: %v", executable, err))
		log.Warn(po.Get("Reversing changes..."))
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
		return errors.New(po.Get("error finishing the progress bar: %s", err.Error()))
	}

	if err = os.Remove(oldExec); err != nil {
		return errors.New(po.Get("error removing old executable: %s", err.Error()))
	}

	if checkos.Unix {
		log.Info(po.Get("Updating completions..."))
		err = updateCompletions(executable)
		if err != nil {
			return errors.New(po.Get("error updating the completions: %s", err.Error()))
		}
	}

	log.Info(po.Get("Upgrade completed successfully"))
	fmt.Printf(
		"%s -> %s\n",
		color.Red.Render(meta.Version),
		color.Green.Render(status.Release.GetTagName()),
	)

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
	log.Warn(po.Get("%s detected, reversing changes before finishing program...", c.String()))
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
