package main

import (
	"embed"

	"github.com/Tom5521/fsize/settings"
	"github.com/jeandeaual/go-locale"
	"github.com/leonelquinteros/gotext"
	"github.com/spf13/viper"
)

//go:embed po
var podir embed.FS

func initLocales() {
	code := viper.GetString(settings.Language)
	if code == "default" {
		var err error
		code, err = locale.GetLanguage()
		if err != nil {
			code = "en"
		}
	}

	loc := gotext.NewLocaleFSWithPath(code, podir, "po")
	loc.AddDomain("default")
	gotext.SetDomain("default")
	gotext.SetLocales([]*gotext.Locale{loc})
	gotext.SetLanguage(code)
}
