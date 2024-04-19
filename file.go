package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/Tom5521/GoNotes/pkg/messages"
)

type File struct {
	Size int64

	Name    string
	AbsPath string
	ModTime time.Time
	IsDir   bool
	Perms   fs.FileMode

	// IsDir vars
	FilesNumber int64
}

func Read(path string) (File, error) {
	var f File

	ofile, err := os.Open(path)
	if err != nil {
		return f, err
	}
	finfo, err := ofile.Stat()
	if err != nil {
		return f, err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return f, err
	}

	f.Size = finfo.Size()
	if finfo.IsDir() && !*NoWalk {
		err = filepath.Walk(path, func(name string, info fs.FileInfo, err error) error {
			if err != nil {
				messages.Warning(err)
				return nil
			}
			if *PrintOnWalk {
				messages.Infof("Reading \"%s\"...", name)
			}

			f.Size += info.Size()
			f.FilesNumber++

			return nil
		})
		if err != nil {
			return f, err
		}
	}

	f.Name = finfo.Name()
	f.IsDir = finfo.IsDir()
	f.ModTime = finfo.ModTime()
	f.Perms = finfo.Mode().Perm()
	f.AbsPath = absPath

	return f, nil
}
