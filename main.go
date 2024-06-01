package main

import (
	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/Tom5521/fsize/filecount"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func main() {
	// Initialize variables
	filecount.Warning = Warning
	filecount.PrintOnWalk = &PrintOnWalk
	err := LoadSettings()
	if err != nil {
		return
	}
	InitFlags()
	root.SetErrPrefix(color.Red.Render("ERROR:"))
	root.Execute()
}

func RunE(cmd *cobra.Command, args []string) (err error) {
	err = cmd.Flags().Parse(args)
	if err != nil {
		return
	}

	switch {
	case GenBashCompletion || GenFishCompletion || GenZshCompletion:
		err = GenerateCompletions(cmd, args)
	case PrintSettingsFlag:
		PrintSettings()
	case len(SettingsFlag) != 0:
		err = ParseSettings(SettingsFlag)
		if err != nil {
			msg.Info("Available configuration keys:")
			PrintSettings()
		}
	default:
		if len(args) == 0 {
			Warning("No file/directory was specified, the current directory will be used. (.)")
			args = append(args, ".")
		}
		for _, f := range args {
			var file File
			file, err = Read(f)
			if err != nil {
				return
			}
			Print(file)
		}
	}

	return
}
