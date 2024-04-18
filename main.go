package main

import (
	"flag"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
)

func main() {
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		msg.Warning("No file/directory was specified, the current directory will be used. (.)")
		files = append(files, ".")
	}
	for _, f := range files {
		file, err := Read(f)
		if err != nil {
			msg.FatalError(err)
		}
		Print(file)
	}
}
