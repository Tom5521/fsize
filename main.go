package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/meta"
	"github.com/Tom5521/fsize/settings"
	"github.com/Tom5521/fsize/stat"
	"github.com/Tom5521/fsize/update"
	"github.com/Tom5521/fsize/walk"
	"github.com/charmbracelet/log"
	"github.com/gookit/color"
	po "github.com/leonelquinteros/gotext"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	releaseTarget string
	upgradable    bool
)

func init() {
	upgradable = releaseTarget == "github-bin"
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize logger.
	cfg := log.Options{
		TimeFormat:      time.Kitchen,
		ReportTimestamp: true,
	}
	if flags.Test {
		cfg.ReportCaller = true
		cfg.Level = log.DebugLevel
	}

	logger := log.NewWithOptions(os.Stderr, cfg)
	log.SetDefault(logger)

	if !flags.Test {
		defer func() {
			if r := recover(); r != nil {
				logger.Fatal(r)
			}
		}()
	}

	// Initialize variables.
	if err := settings.InitSettings(); err != nil {
		log.Fatal(err)
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
		settings.Print()
	case len(flags.SettingsFlag) != 0:
		err = settings.Parse(flags.SettingsFlag)
		if err != nil {
			log.Info(po.Get("Available configuration keys:"))
			settings.Print()
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
			log.Info(po.Get("Already in latest version"))
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
		fmt.Println(po.Get("Self-upgradable:"), upgradable)
		fmt.Println(po.Get("Version:"), meta.LongVersion)
		fmt.Println(po.Get("Updated:"), updated)
		if !updated {
			fmt.Println(po.Get("Latest version:"), tag)
		}

		fmt.Println(po.Get("Source Code: %s", "https://github.com/Tom5521/fsize"))
	default:
		if len(args) == 0 {
			log.Warn(
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
			walk.ProcessFile(&file)

			fmt.Println(file)
		}
	}

	return err
}
