package main

import (
	"flag"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
)

func main() {
	err := LoadSettings()
	if err != nil {
		msg.Error(err)
	}
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
}
