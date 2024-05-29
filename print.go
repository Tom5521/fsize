package main

import (
	"fmt"
	"time"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/Tom5521/fsize/checkos"
	"github.com/gookit/color"
	cbytes "github.com/labstack/gommon/bytes"
)

func Print(f File) {
	makePrint("Name:", f.Name)
	makePrint("Size:", cbytes.New().Format(f.Size))
	makePrint("Absolute Path:", f.AbsPath)
	makePrint("Date Modified:", f.ModTime.Format(time.DateTime))
	makePrint("Is directory:", f.IsDir)
	makePrint("Permissions:", fmt.Sprintf("%v/%v", int(f.Perms), f.Perms))
	if f.IsDir && !NoWalk {
		makePrint("Number of files:", f.FilesNumber)
	}

	switch {
	case checkos.Unix:
		makePrint("UID/Name:", fmt.Sprintf("%v/%v", f.User.Uid, f.User.Username))
		makePrint("GID/Name:", fmt.Sprintf("%v/%v", f.Group.Gid, f.Group.Name))
	case checkos.Windows:
		makePrint("Creation Date:", f.CreationDate.Format(time.DateTime))
	}
}

func PrintSettings() {
	printSetting := func(key string) {
		s := Settings.Read(key)
		_, isBool := s.(bool)
		if isBool {
			switch s {
			case true:
				s = color.Green.Render(s)
			default:
				s = color.Red.Render(s)
			}
		}
		fmt.Print(key + ": ")
		fmt.Println(s)
	}
	for _, s := range Keys {
		printSetting(s)
	}
}

func Warning(warn ...any) {
	if NoWarns {
		return
	}
	msg.Warning(warn...)
}

func makePrint(title string, content ...any) {
	color.Green.Print(title + " ")
	fmt.Println(content...)
}
