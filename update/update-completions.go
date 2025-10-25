package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	po "github.com/leonelquinteros/gotext"
)

var completionPaths = map[string]string{
	"bash": fmt.Sprintf("%s/share/bash-completion/completions/fsize", xdg.DataHome),
	"fish": fmt.Sprintf("%s/share/fish/vendor_completions.d/fsize.fish", xdg.DataHome),
	"zsh":  fmt.Sprintf("%s/share/zsh/site-functions/_fsize", xdg.DataHome),
}

func removeCompletions() {
	for k, v := range completionPaths {
		if _, err := exec.LookPath(k); err != nil {
			continue
		}
		err := os.Remove(v)
		if err != nil {
			log.Warn(po.Get("error removing %s completions: %s", k, err.Error()))
		}
	}
}

func updateCompletions(executable string) (err error) {
	var confirm bool
	err = huh.NewConfirm().
		Title(po.Get("Do you want to update the completions?")).
		Affirmative(po.Get("Yes")).
		Negative(po.Get("No")).
		Value(&confirm).
		WithTheme(huh.ThemeBase()).
		RunAccessible(os.Stderr, os.Stdin)
	if !confirm || err != nil {
		if err != nil {
			log.Error(po.Get("error running confirm dialog"), "err", err)
		}
		return err
	}
	removeCompletions()
	for shell, path := range completionPaths {
		_, exists := exec.LookPath(shell)
		if exists != nil {
			log.Warn(po.Get("%s not found.", shell))
			continue
		}
		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			log.Warn(po.Get("error creating completion directory: %s", err.Error()))
			continue
		}
		cmd := exec.Command(executable, fmt.Sprintf("--gen-%s-completion", shell), path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return err
}
