package update

import (
	"os/exec"
	"strings"

	"github.com/Tom5521/fsize/echo"
)

func updateCompletions() (err error) {
	const instructions string = `bash|--gen-bash-completion /usr/share/bash-completion/completions/fsize
fish|--gen-fish-completion /usr/share/fish/vendor_completions.d/fsize.fish
zsh|--gen-zsh-completion /usr/share/zsh/site-functions/_fsize`

	lines := strings.Split(instructions, "\n")
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
