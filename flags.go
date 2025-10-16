package main

import (
	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/meta"
	"github.com/Tom5521/fsize/settings"
	"github.com/charmbracelet/log"
	"github.com/gookit/color"
	"github.com/leonelquinteros/gotext"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var root cobra.Command

func initRoot() {
	root = cobra.Command{
		Use:   "fsize",
		Short: gotext.Get("Displays the file/folder properties."),

		PreRunE: func(cmd *cobra.Command, args []string) error {
			lvl, err := log.ParseLevel(flags.LogLevel)
			if err != nil {
				return err
			}
			log.SetLevel(lvl)

			if flags.NoWarns {
				log.SetLevel(log.ErrorLevel)
			}
			return nil
		},
		RunE: RunE,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return viper.WriteConfig()
		},
		Version: meta.LongVersion,
	}
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
	flag.StringSliceVar(&flags.SettingsFlag, "config", []string{},
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
	if upgradable {
		flag.BoolVar(
			&flags.Update,
			"update",
			false,
			gotext.Get(
				`Automatically updates the program by overwriting the binary and regenerating the completions.`,
			),
		)
	}
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
	flag.BoolVarP(
		&flags.NoProgress,
		"no-progress",
		"s",
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
	flag.StringSliceVarP(
		&flags.Patterns,
		"pattern", "f",
		viper.GetStringSlice(settings.Pattern),
		gotext.Get(
			`If the pattern is not "", only files that match it will be included in the count.
The pattern must be a regular expression unless the --wildcard flag on`,
		),
	)

	flag.StringSliceVarP(
		&flags.IgnorePatterns,
		"ignore", "i",
		viper.GetStringSlice(settings.IgnorePattern),
		gotext.Get(`If ignore is not "", the files that match it will be excluded from the count.
The pattern must be a regular expression unless the --wildcard flag is on`),
	)
	flag.BoolVarP(
		&flags.Wildcard,
		"wildcard",
		"w",
		viper.GetBool(settings.Wildcard),
		gotext.Get(
			`Switches --ignore & --pattern from regular expressions to wildcard patterns`,
		),
	)
	flag.BoolVarP(
		&flags.NotClearBar,
		"not-clear-bar", "c",
		viper.GetBool(settings.NotClearBar),
		gotext.Get(`Prevents the progress indicator from being cleared after finishing`),
	)
	flag.UintVarP(
		&flags.Depth,
		"depth",
		"d",
		viper.GetUint(settings.Depht),
		gotext.Get(
			`Indicates the maximum depth to traverse within a directory; 
files/directories deeper than this will be skipped`,
		),
	)
	flag.IntVarP(
		&flags.WarnLimit,
		"warn-limit",
		"x",
		viper.GetInt(settings.WarnLimit),
		gotext.Get(
			`Indicates the maximum number of warnings that will be printed. 
If it is -1, there is no limit.`,
		),
	)
	flag.StringVar(
		&flags.LogLevel,
		"log",
		viper.GetString(settings.LogLevel),
		gotext.Get(
			`Indicates the log level, which can be debug, info, warn, error, or fatal.`,
		),
	)
	flag.BoolVarP(&flags.FollowSymlinks, "follow-symlinks", "l", false,
		gotext.Get(`If enabled, the program will follow symbolic links`),
	)

	root.MarkFlagsMutuallyExclusive(
		"progress",
		"no-progress",
	)
	root.MarkFlagsMutuallyExclusive(
		"no-progress",
		"progress-delay",
		"print-on-walk",
	)

	otherFlags := []string{
		"gen-zsh-completion",
		"gen-bash-completion",
		"gen-fish-completion",
		"print-settings",
		"bin-info",
		"config",
	}
	if upgradable {
		otherFlags = append(otherFlags, "update")
	}
	root.MarkFlagsMutuallyExclusive(otherFlags...)
}
