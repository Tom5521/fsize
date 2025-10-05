package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Tom5521/fsize/echo"
	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/meta"
	"github.com/Tom5521/fsize/settings"
	"github.com/Tom5521/fsize/stat"
	"github.com/Tom5521/fsize/update"
	"github.com/gookit/color"
	po "github.com/leonelquinteros/gotext"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	defer func() {
		if r := recover(); r != nil {
			color.Errorln(r)
		}
	}()

	// Initialize variables

	if err := settings.InitSettings(); err != nil {
		echo.Error(err)
		return
	}

	initLocales()
	initRoot()

	InitFlags()
	root.SetErrPrefix(color.Error.Render(po.Get("ERROR:")))
	defer root.Execute()
}

func RunE(cmd *cobra.Command, args []string) (err error) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		color.Enable = false
	}

	switch {
	case flags.GenBashCompletion || flags.GenFishCompletion || flags.GenZshCompletion:
		err = GenerateCompletions(cmd, args)
	case flags.PrintSettingsFlag:
		echo.Settings()
	case len(flags.SettingsFlag) != 0:
		err = settings.Parse(flags.SettingsFlag)
		if err != nil {
			echo.Info(po.Get("Available configuration keys:"))
			echo.Settings()
		}
	case flags.Update:
		var (
			tag     string
			updated bool
		)
		tag, updated, err = update.CheckUpdate()
		if err != nil {
			return err
		}
		if updated {
			echo.Info(po.Get("Already in latest version"))
			return err
		}
		err = update.ApplyUpdate(tag)
	case flags.BinInfo:
		var (
			tag     string
			updated bool
		)
		tag, updated, err = update.CheckUpdate()
		if err != nil {
			return err
		}

		fmt.Println("GOOS:", runtime.GOOS)
		fmt.Println("GOARCH:", runtime.GOARCH)
		fmt.Println(po.Get("Version:"), meta.LongVersion)
		fmt.Println(po.Get("Updated:"), updated)
		if !updated {
			fmt.Println(po.Get("Latest version:"), tag)
		}

		fmt.Println(po.Get("Source Code: %s", "https://github.com/Tom5521/fsize"))
	default:
		if len(args) == 0 {
			echo.Warning(
				po.Get(
					"No file/directory was specified, the current directory will be used. (.)",
				),
			)
			args = append(args, ".")
		}
		for _, f := range args {
			var file stat.File
			file, err = stat.NewFile(f)
			if err != nil {
				return err
			}
			fmt.Println(file)
		}
	}

	return err
}
