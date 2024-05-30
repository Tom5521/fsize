package main

import (
	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func main() {
	err := LoadSettings()
	if err != nil {
		return
	}
	InitFlags()
	root.SetErrPrefix(color.Red.Render("ERROR:"))
	root.Execute()
}

func RunE(cmd *cobra.Command, args []string) (err error) {
	err = cmd.PersistentFlags().Parse(args)
	if err != nil {
		return err
	}
	if len(SettingsFlag) != 0 {
		err = ParseSettings(SettingsFlag)
		if err != nil {
			msg.Info("Available configuration keys:")
			PrintSettings()
		}
		return err
	}
	if PrintSettingsFlag {
		PrintSettings()
		return nil
	}

	if len(args) == 0 {
		Warning("No file/directory was specified, the current directory will be used. (.)")
		args = append(args, ".")
	}
	for _, f := range args {
		var file File
		file, err = Read(f)
		if err != nil {
			return err
		}
		Print(file)
	}
	return nil
}
