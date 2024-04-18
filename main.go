package main

import (
	"os"

	"github.com/Tom5521/GoNotes/pkg/messages"
)

func main() {
	var file File
	var err error
	if len(os.Args) == 1 {
		messages.Warning("No file was specified, the current directory will be used. (.)")
		file, err = Read(".")
	}
	file, err = Read(os.Args[1])
	if err != nil {
		messages.FatalError(err)
	}
	Print(file)
}
