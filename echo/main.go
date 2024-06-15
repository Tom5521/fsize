package echo

import (
	"fmt"

	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/settings"
	conf "github.com/Tom5521/goconf"
	"github.com/gookit/color"
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
	color.Warnln(warn...)
}
