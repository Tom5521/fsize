package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Tom5521/fsize/echo"
)

func updateCompletions(executable string) (err error) {
	instructions := map[string]string{
		"bash": "/usr/local/share/bash-completion/completions/fsize",
		"fish": "/usr/local/share/fish/vendor_completions.d/fsize.fish",
		"zsh":  "/usr/local/share/zsh/site-functions/_fsize",
	}

	if isMaybeRunningInTermux() {
		for k, v := range instructions {
			instructions[k] = filepath.Join(os.Getenv("PREFIX"), v)
		}
	}

	for shell, path := range instructions {
		_, exists := exec.LookPath(shell)
		if exists != nil {
			echo.Warningf("%s not found.", shell)
			continue
		}
		cmd := exec.Command(executable, fmt.Sprintf("--gen-%s-completion", shell), path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return
		}
	}

	return
}
