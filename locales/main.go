package locales

import (
	"embed"
	"errors"
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

	if code == "en" || err != nil {
		Po.ParseFile("po/en.pot")
	} else {
		file := fmt.Sprintf("po/%s.po", code)

		if _, err := fs.Stat(podir, file); errors.Is(err, fs.ErrNotExist) {
			file = "po/en.pot"
		}

		Po.ParseFile(file)
	}
}
