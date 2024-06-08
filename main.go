package main

import (
	"fmt"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/Tom5521/fsize/echo"
	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/meta"
	"github.com/Tom5521/fsize/settings"
	"github.com/Tom5521/fsize/stat"
	"github.com/Tom5521/fsize/update"
	conf "github.com/Tom5521/goconf"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var Settings conf.Preferences

func main() {
	// Initialize variables
	err := settings.Load()
	if err != nil {
		return
	}
	Settings = settings.Settings
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
	case flags.GenBashCompletion || flags.GenFishCompletion || flags.GenZshCompletion:
		err = GenerateCompletions(cmd, args)
	case flags.PrintSettingsFlag:
		echo.Settings(Settings)
	case len(flags.SettingsFlag) != 0:
		err = settings.Parse(flags.SettingsFlag)
		if err != nil {
			msg.Info("Available configuration keys:")
			echo.Settings(Settings)
		}
	case flags.Update:
		var (
			tag     string
			updated bool
		)
		tag, updated, err = update.CheckUpdate()
		if err != nil {
			return
		}
		if updated {
			msg.Info("Already in latest version")
			return
		}
		err = update.ApplyUpdate(tag)
	case flags.BinInfo:
		var (
			tag     string
			updated bool
		)
		tag, updated, err = update.CheckUpdate()
		if err != nil {
			return
		}

		fmt.Println("Version:", meta.Version)
		fmt.Println("Updated:", updated)
		if !updated {
			fmt.Println("Latest version:", tag)
		}
	default:
		if len(args) == 0 {
			echo.Warning("No file/directory was specified, the current directory will be used. (.)")
			args = append(args, ".")
		}
		for _, f := range args {
			var file stat.File
			file, err = stat.NewFile(f)
			if err != nil {
				return
			}
			echo.File(file)
		}
	}

	return
}
