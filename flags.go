package main

import "flag"

var (
	PrintOnWalk = flag.Bool("print-on-walk", false,
		"Prints the name of the file being walked if a directory has been selected.",
	)
	NoWalk = flag.Bool("no-walk", false,
		"Skips walking inside the directories.",
	)
)
