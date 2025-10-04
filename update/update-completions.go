package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Tom5521/fsize/echo"
	"github.com/adrg/xdg"
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
			echo.Warning(po.Get("error removing %s completions: %s", k, err.Error()))
		}
	}
}

func updateCompletions(executable string) (err error) {
	var confirm bool
	prompt := &survey.Confirm{
		Message: "Do you want to update the completions?",
	}
	survey.AskOne(prompt, &confirm)
	if !confirm {
		return nil
	}
	removeCompletions()
	for shell, path := range completionPaths {
		_, exists := exec.LookPath(shell)
		if exists != nil {
			echo.Warning(po.Get("%s not found.", shell))
			continue
		}
		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			echo.Warning(po.Get("error creating completion directory: %s", err.Error()))
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
