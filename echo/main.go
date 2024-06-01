package echo

import (
	"fmt"
	"time"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/Tom5521/fsize/checkos"
	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/settings"
	"github.com/Tom5521/fsize/stat"
	conf "github.com/Tom5521/goconf"
	"github.com/gookit/color"
	cbytes "github.com/labstack/gommon/bytes"
)

func Settings(s conf.Preferences) {
	printSetting := func(key string) {
		s := s.Read(key)
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
	for _, s := range settings.Keys {
		printSetting(s)
	}
}

func Warning(warn ...any) {
	if flags.NoWarns {
		return
	}
	msg.Warning(warn...)
}

func Make(title string, content ...any) {
	color.Green.Print(title + " ")
	fmt.Println(content...)
}

func File(f stat.File) {
	Make("Name:", f.Name)
	Make("Size:", cbytes.New().Format(f.Size))
	Make("Absolute Path:", f.AbsPath)
	Make("Date Modified:", f.ModTime.Format(time.DateTime))
	Make("Is directory:", f.IsDir)
	Make("Permissions:", fmt.Sprintf("%v/%v", int(f.Perms), f.Perms))
	if f.IsDir && !flags.NoWalk {
		Make("Number of files:", f.FilesNumber)
	}

	switch {
	case checkos.Unix:
		Make("UID/Name:", fmt.Sprintf("%v/%v", f.User.Uid, f.User.Username))
		Make("GID/Name:", fmt.Sprintf("%v/%v", f.Group.Gid, f.Group.Name))
	case checkos.Windows:
		Make("Creation Date:", f.CreationDate.Format(time.DateTime))
	}
}
