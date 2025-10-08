package flags

import "time"

var (
	PrintOnWalk    bool
	NoWalk         bool
	Progress       bool
	NoWarns        bool
	Update         bool
	BinInfo        bool
	NoProgress     bool
	ProgressDelay  time.Duration
	Wildcard       bool
	Patterns       []string
	IgnorePatterns []string

	// Hidden flags.

	Test bool

	// Shell completions.

	GenBashCompletion bool
	GenFishCompletion bool
	GenZshCompletion  bool

	PrintSettingsFlag bool
	SettingsFlag      []string
)
