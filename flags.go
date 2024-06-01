package main

import (
	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/meta"
	"github.com/Tom5521/fsize/settings"
	"github.com/spf13/cobra"
)

var root = cobra.Command{
	Use:     "fsize",
	Short:   "Displays the file/folder properties.",
	RunE:    RunE,
	Version: meta.Version,
}

func InitFlags() {
	flag := root.Flags()
	flag.BoolVar(&flags.PrintOnWalk, "print-on-walk", Settings.Bool(settings.AlwaysPrintOnWalk),
		"Prints the name of the file being walked if a directory has been selected.",
	)
	flag.BoolVar(&flags.NoWalk, "no-walk", Settings.Bool(settings.AlwaysSkipWalk),
		"Skips walking inside the directories.",
	)
	flag.BoolVarP(&flags.Progress, "progress", "p", Settings.Bool(settings.AlwaysShowProgress),
		"Displays a file count and progress bar when counting and summing file sizes.",
	)
	flag.StringSliceVarP(&flags.SettingsFlag, "config", "c", []string{},
		`Configure the variables used for preferences
Example: "fsize --config 'AlwaysShowProgress=true,AlwaysPrintOnWalk=false'".

To see the available variables and their values run "fsize --print-settings".`,
	)
	flag.BoolVar(&flags.PrintSettingsFlag, "print-settings", false,
		"Prints the current configuration values.",
	)
	flag.BoolVar(&flags.NoWarns, "no-warns", Settings.Bool(settings.HideWarnings),
		"Hide possible warnings.",
	)
	flag.BoolVar(&flags.GenBashCompletion, "gen-bash-completion", false,
		`Generate a completion file for bash
if any, the first argument will be taken as output file.`,
	)
	flag.BoolVar(&flags.GenFishCompletion, "gen-fish-completion", false,
		`Generate a completion file for fish
if any, the first argument will be taken as output file.`,
	)
	flag.BoolVar(&flags.GenZshCompletion, "gen-zsh-completion", false,
		`Generate a completion file for zsh
if any, the first argument will be taken as output file.`,
	)
}
