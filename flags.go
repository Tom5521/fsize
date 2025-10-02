package main

import (
	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/meta"
	"github.com/Tom5521/fsize/settings"
	"github.com/gookit/color"
	"github.com/leonelquinteros/gotext"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var root = cobra.Command{
	Use:   "fsize",
	Short: gotext.Get("Displays the file/folder properties."),
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return viper.WriteConfig()
	},
	RunE:    RunE,
	Version: meta.LongVersion,
}

func InitFlags() {
	flag := root.Flags()

	flag.BoolVar(&flags.PrintOnWalk, "print-on-walk", viper.GetBool(settings.AlwaysPrintOnWalk),
		gotext.Get("Prints the name of the file being walked if a directory has been selected."),
	)
	flag.BoolVar(&flags.NoWalk, "no-walk", viper.GetBool(settings.AlwaysSkipWalk),
		gotext.Get("Skips walking inside the directories."),
	)
	flag.BoolVarP(&flags.Progress, "progress", "p", viper.GetBool(settings.AlwaysShowProgress),
		gotext.Get("Displays a file count and progress bar when counting and summing file sizes."),
	)
	flag.StringSliceVarP(&flags.SettingsFlag, "config", "c", []string{},
		gotext.Get(`Configure the variables used for preferences
Example: "fsize --config 'always-show-progress=true,always-print-on-walk=false'".

To see the available variables and their values run "fsize --print-settings".`,
		))
	flag.BoolVar(&flags.PrintSettingsFlag, "print-settings", false,
		gotext.Get("Prints the current configuration values."),
	)
	flag.BoolVar(&flags.NoWarns, "no-warns", viper.GetBool(settings.HideWarnings),
		gotext.Get("Hide possible warnings."),
	)
	flag.BoolVar(&flags.GenBashCompletion, "gen-bash-completion", false,
		gotext.Get(`Generate a completion file for bash
if any, the first argument will be taken as output file.`),
	)
	flag.BoolVar(&flags.GenFishCompletion, "gen-fish-completion", false,
		gotext.Get(`Generate a completion file for fish
if any, the first argument will be taken as output file.`),
	)
	flag.BoolVar(&flags.GenZshCompletion, "gen-zsh-completion", false,
		gotext.Get(`Generate a completion file for zsh
if any, the first argument will be taken as output file.`),
	)
	flag.BoolVar(
		&flags.Update,
		"update",
		false,
		gotext.Get(
			`Automatically updates the program by overwriting the binary and regenerating the completions.`,
		),
	)
	flag.BoolVar(
		&flags.BinInfo,
		"bin-info",
		false,
		gotext.Get("Displays the information of the binary"),
	)
	flag.BoolVar(
		&color.Enable,
		"color",
		!viper.GetBool(settings.NoColor),
		"enable or disable the color",
	)
	flag.BoolVar(
		&flags.NoProgress,
		"no-progress",
		viper.GetBool(settings.HideProgress),
		gotext.Get("Disable any progress indicator."),
	)

	flag.DurationVar(
		&flags.ProgressDelay,
		"progress-delay",
		viper.GetDuration(settings.ProgressDelay),
		gotext.Get(`Specifies how long the program should be counting files
before a progress indicator appears`),
	)

	flag.BoolVar(&flags.Test, "test", false, "---")
	flag.MarkHidden("test")

	root.MarkFlagsMutuallyExclusive(
		"update",
		"gen-zsh-completion",
		"gen-bash-completion",
		"gen-fish-completion",
		"print-settings",
		"bin-info",
		"config",
	)
}
