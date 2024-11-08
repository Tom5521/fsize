package echo

import (
	"fmt"

	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/locales"
	"github.com/Tom5521/fsize/settings"
	conf "github.com/Tom5521/goconf"
	"github.com/gookit/color"
)

var po = locales.Po

func Settings(s *conf.Preferences) {
	printSetting := func(key string) {
		s := s.Get(key)
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

func Println(txt string, args ...any) {
	fmt.Print(po.Get(txt))
	if len(args) > 0 {
		fmt.Print(args...)
	}
	fmt.Println()
}

func Printfln(format string, a ...any) {
	fmt.Println(po.Get(format, a...))
}

func Warningf(format string, arg ...any) {
	Warning(po.Get(format, arg...))
}

func Warning(warn string) {
	if flags.NoWarns {
		return
	}
	color.Warnln(po.Get(warn))
}

func Info(format string, a ...any) {
	color.Info.Println(po.Get(format, a...))
}

func Error(format string, a ...any) {
	color.Error.Println(po.Get(format, a...))
}
