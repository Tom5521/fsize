package main

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Tom5521/fsize/settings"
	"github.com/jeandeaual/go-locale"
	"github.com/leonelquinteros/gotext"
	"github.com/spf13/viper"
)

//go:embed po
var podir embed.FS

func loadLocalesToTmp(tmpDir string) error {
	return fs.WalkDir(podir, "po",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return os.MkdirAll(filepath.Join(tmpDir, path), os.ModePerm)
			}
			file, err := podir.ReadFile(path)
			if err != nil {
				return err
			}

			return os.WriteFile(filepath.Join(tmpDir, path), file, os.ModePerm)
		})
}

func initLocales() {
	var err error
	code := viper.GetString(settings.Language)
	if code == "default" {
		code, err = locale.GetLanguage()
		if err != nil {
			code = "en"
		}
	}
	tmpDir, err := os.MkdirTemp(os.TempDir(), "fsize-locales")
	if err != nil {
		return
	}
	err = loadLocalesToTmp(tmpDir)
	if err != nil {
		return
	}
	defer os.RemoveAll(tmpDir)

	gotext.Configure(filepath.Join(tmpDir, "po"), code, "default")
}
