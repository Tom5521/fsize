package update

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Tom5521/fsize/echo"
)

func updateCompletions() (err error) {
	const instructions = `bash|/usr/local/share/bash-completion/completions/fsize
fish|/usr/local/share/fish/vendor_completions.d/fsize.fish
zsh|/usr/local/share/zsh/site-functions/_fsize`

	for _, line := range strings.Split(instructions, "\n") {
		parts := strings.SplitN(line, "|", 2)
		shell, path := parts[0], parts[1]

		_, exists := exec.LookPath(shell)
		if exists != nil {
			echo.Warning(shell, po.Get("not found."))
			continue
		}
		cmd := exec.Command("fsize", fmt.Sprintf("--gen-%s-completion", shell), path)
		err = cmd.Run()
		if err != nil {
			return
		}
	}

	return
}
