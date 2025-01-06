package main

import (
	"github.com/Tom5521/fsize/flags"
	"github.com/spf13/cobra"
)

func GenerateCompletions(cmd *cobra.Command, args []string) (err error) {
	var (
		customName     bool
		filename       string
		completionFunc func(string) error
	)
	if len(args) == 0 {
		customName = false
		filename = "fsize-completions"
	} else {
		customName = true
		filename = args[0]
	}

	switch {
	case flags.GenZshCompletion:
		if !customName {
			filename += ".zsh"
		}
		completionFunc = cmd.GenZshCompletionFile
	case flags.GenFishCompletion:
		if !customName {
			filename += ".fish"
		}
		completionFunc = func(filename string) error {
			return cmd.GenFishCompletionFile(filename, true)
		}
	case flags.GenBashCompletion:
		if !customName {
			filename += ".sh"
		}
		completionFunc = func(s string) error {
			return cmd.GenBashCompletionFileV2(s, true)
		}
	}

	err = completionFunc(filename)

	return
}
