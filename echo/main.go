package echo

import (
	"fmt"
	"os"

	"github.com/Tom5521/fsize/flags"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

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

func Warningf(format string, arg ...any) {
	Warning(fmt.Sprintf(format, arg...))
}

func Warning(warn string) {
	if flags.NoWarns {
		return
	}
	warn = color.Warn.Render(warn)
	fmt.Fprintln(os.Stderr, warn)
}

func Info(a ...any) {
	info := color.Info.Render(a...)
	fmt.Fprintln(os.Stderr, info)
}

func Error(a ...any) {
	err := color.Error.Render(a...)
	fmt.Fprintln(os.Stderr, err)
}
