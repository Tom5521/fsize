package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/schollz/progressbar/v3"
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

	if *ProgressBar && finfo.IsDir() && !*NoWalk {
		msg.Info("Counting the amount of files...")
		var fileNumber int64
		var files []string
		err = filepath.Walk(path, func(name string, _ fs.FileInfo, err error) error {
			if err != nil {
				msg.Warning(err)
				return nil
			}
			fileNumber++
			files = append(files, name)
			fmt.Printf("%v files found...\r", fileNumber)
			return nil
		})
		if err != nil {
			return f, err
		}
		msg.Info("Calculating total size...")
		f.FilesNumber = fileNumber
		bar := progressbar.Default(fileNumber)
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				bar.Add(1)
				continue
			}
			f.Size += info.Size()
			bar.Add(1)
		}
	} else if finfo.IsDir() && !*NoWalk {
		err = filepath.Walk(path, func(name string, info fs.FileInfo, err error) error {
			if err != nil {
				msg.Warning(err)
				return nil
			}
			if *PrintOnWalk {
				msg.Infof("Reading \"%s\"...", name)
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
