package main

import (
	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/spf13/cobra"
)

func main() {
	err := LoadSettings()
	if err != nil {
		msg.Error(err)
	}
	InitFlags()
	root.RunE = func(cmd *cobra.Command, args []string) error {
		root.PersistentFlags().Parse(args)
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
			file, err := Read(f)
			if err != nil {
				return err
			}
			Print(file)
		}
		return nil
	}
	err = root.Execute()
	if err != nil {
		msg.FatalError(err)
	}
	/*
		InitFlags()
		if SettingsFlag != "" {
			err = ParseSettings(SettingsFlag)
			if err != nil {
				msg.Error(err)
				msg.Info("Available configuration keys:")
				PrintSettings()
			}
			return
		}
		if PrintSettingsFlag {
			PrintSettings()
			return
		}
		files := flag.Args()
		if len(files) == 0 {
			Warning("No file/directory was specified, the current directory will be used. (.)")
			files = append(files, ".")
		}
		for _, f := range files {
			file, err := Read(f)
			if err != nil {
				msg.FatalError(err)
			}
			Print(file)
		}
	*/
}
