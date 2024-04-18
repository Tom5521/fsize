package main

import "flag"

var (
	PrintOnWalk = flag.Bool("print-on-walk", false, "")
	NoWalk      = flag.Bool("no-walk", false, "")
)
