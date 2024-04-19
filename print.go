package main

import (
	"fmt"
	"time"

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
	if f.IsDir && !*NoWalk {
		makePrint("Number of files:", f.FilesNumber)
	}
}

func makePrint(title string, content ...any) {
	color.Green.Print(title + " ")
	fmt.Println(content...)
}
