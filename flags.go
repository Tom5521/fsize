package main

import (
	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/meta"
	"github.com/Tom5521/fsize/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var root = cobra.Command{
	Use:   "fsize",
	Short: po.Get("Displays the file/folder properties."),
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return viper.WriteConfig()
	},
	RunE:    RunE,
	Version: meta.LongVersion,
}

func InitFlags() {
	flag := root.Flags()

	flag.BoolVar(&flags.PrintOnWalk, "print-on-walk", viper.GetBool(settings.AlwaysPrintOnWalk),
		po.Get("Prints the name of the file being walked if a directory has been selected."),
	)
	flag.BoolVar(&flags.NoWalk, "no-walk", viper.GetBool(settings.AlwaysSkipWalk),
		po.Get("Skips walking inside the directories."),
	)
	flag.BoolVarP(&flags.Progress, "progress", "p", viper.GetBool(settings.AlwaysShowProgress),
		po.Get("Displays a file count and progress bar when counting and summing file sizes."),
	)
	flag.StringSliceVarP(&flags.SettingsFlag, "config", "c", []string{},
		po.Get(`Configure the variables used for preferences
Example: "fsize --config 'AlwaysShowProgress=true,AlwaysPrintOnWalk=false'".

To see the available variables and their values run "fsize --print-settings".`,
		))
	flag.BoolVar(&flags.PrintSettingsFlag, "print-settings", false,
		po.Get("Prints the current configuration values."),
	)
	flag.BoolVar(&flags.NoWarns, "no-warns", viper.GetBool(settings.HideWarnings),
		po.Get("Hide possible warnings."),
	)
	flag.BoolVar(&flags.GenBashCompletion, "gen-bash-completion", false,
		po.Get(`Generate a completion file for bash
if any, the first argument will be taken as output file.`),
	)
	flag.BoolVar(&flags.GenFishCompletion, "gen-fish-completion", false,
		po.Get(`Generate a completion file for fish
if any, the first argument will be taken as output file.`),
	)
	flag.BoolVar(&flags.GenZshCompletion, "gen-zsh-completion", false,
		po.Get(`Generate a completion file for zsh
if any, the first argument will be taken as output file.`),
	)
	flag.BoolVar(
		&flags.Update,
		"update",
		false,
		po.Get(
			`Automatically updates the program by overwriting the binary and regenerating the completions.`,
		),
	)
	flag.BoolVar(
		&flags.BinInfo,
		"bin-info",
		false,
		po.Get("Displays the information of the binary"),
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
