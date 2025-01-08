package echo

import (
	"fmt"
	"os"

	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/locales"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

var po = locales.Po

func Settings() {
	for _, key := range viper.AllKeys() {
		s := viper.Get(key)
		switch v := s.(type) {
		case bool:
			if v {
				s = color.Green.Render(s)
			} else {
				s = color.Red.Render(s)
			}
		}
		fmt.Print(key + ": ")
		fmt.Println(s)
	}
}

func Println(txt string, args ...any) {
	fmt.Print(po.Get(txt))
	if len(args) > 0 {
		fmt.Print(" ")
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
	warn = color.Warn.Render(po.Get(warn))
	fmt.Fprintln(os.Stderr, warn)
}

func Info(format string, a ...any) {
	color.Info.Println(po.Get(format, a...))
}

func Error(format string, a ...any) {
	err := color.Error.Render(po.Get(format, a...))
	fmt.Fprintln(os.Stderr, err)
}
