package locales

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/jeandeaual/go-locale"
	"github.com/leonelquinteros/gotext"
)

var (
	//go:embed po
	podir embed.FS
	Po    = gotext.NewPoFS(podir)
)

func init() {
	code, err := locale.GetLanguage()

	file := "po/en.pot" // Default language.

	if code != "en" && err == nil {
		langFile := fmt.Sprintf("po/%s.po", code)

		// Check if the language exists.
		if _, err = fs.Stat(podir, langFile); err == nil {
			file = langFile
		}
	}

	Po.ParseFile(file)
}
