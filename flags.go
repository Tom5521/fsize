package main

import (
	"github.com/Tom5521/fsize/meta"
	"github.com/spf13/cobra"
)

var (
	root = cobra.Command{
		Use:     "fsize",
		Short:   "Displays the file/folder properties.",
		RunE:    RunE,
		Version: meta.Version,
	}

	PrintOnWalk bool
	NoWalk      bool
	Progress    bool
	NoWarns     bool

	PrintSettingsFlag bool
	SettingsFlag      []string
)

func InitFlags() {
	flag := root.Flags()
	flag.BoolVar(&PrintOnWalk, "print-on-walk", Settings.Bool(AlwaysPrintOnWalk),
		"Prints the name of the file being walked if a directory has been selected.",
	)
	flag.BoolVar(&NoWalk, "no-walk", Settings.Bool(AlwaysSkipWalk),
		"Skips walking inside the directories.",
	)
	flag.BoolVarP(&Progress, "progress", "p", Settings.Bool(AlwaysShowProgress),
		"Displays a file count and progress bar when counting and summing file sizes.",
	)
	flag.StringSliceVarP(&SettingsFlag, "config", "c", []string{},
		`Configure the variables used for preferences
		Example: "fsize --config 'AlwaysShowProgress=true,AlwaysPrintOnWalk=false'".

		To see the available variables and their values run "fsize --print-settings".`,
	)
	flag.BoolVar(&PrintSettingsFlag, "print-settings", false,
		"Prints the current configuration values.",
	)
	flag.BoolVar(&NoWarns, "no-warns", Settings.Bool(HideWarnings),
		"Hide possible warnings.",
	)
}
