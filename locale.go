package main

import (
	"embed"

	"github.com/jeandeaual/go-locale"
	"github.com/leonelquinteros/gotext"
)

//go:embed po
var podir embed.FS

func init() {
	code, err := locale.GetLanguage()
	if err != nil {
		code = "en"
	}

	loc := gotext.NewLocaleFSWithPath(code, podir, "po")
	loc.AddDomain("default")
	gotext.SetDomain("default")
	gotext.SetLocales([]*gotext.Locale{loc})
	gotext.SetLanguage(code)
}
