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
	NotClearBar    bool
	FollowSymlinks bool
	Wildcard       bool
	ProgressDelay  time.Duration
	Patterns       []string
	IgnorePatterns []string
	LogLevel       string
	Depth          uint
	WarnLimit      int

	// Hidden flags.

	Test bool

	// Shell completions.

	GenBashCompletion bool
	GenFishCompletion bool
	GenZshCompletion  bool

	PrintSettingsFlag bool
	SettingsFlag      []string
)
