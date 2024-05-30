package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"time"

	msg "github.com/Tom5521/GoNotes/pkg/messages"
	"github.com/Tom5521/fsize/stat"
	"github.com/schollz/progressbar/v3"
)

type File struct {
	Size int64

	Name    string
	AbsPath string
	ModTime time.Time
	IsDir   bool
	Perms   fs.FileMode

	Group *user.Group
	User  *user.User

	CreationDate time.Time

	// IsDir vars
	FilesNumber int64
}

func Read(path string) (f File, err error) {
	var (
		finfo   os.FileInfo
		absPath string
	)
	_, finfo, absPath, err = BasicFileInfo(path)
	if err != nil {
		return
	}
	f, err = BasicFile(finfo, absPath)
	if err != nil {
		return
	}

	if Progress && finfo.IsDir() && !NoWalk {
		msg.Info("Counting the amount of files...")
		var (
			fileNumber int64
			files      []os.FileInfo
		)
		err = filepath.Walk(path, func(_ string, info fs.FileInfo, err error) error {
			if err != nil {
				Warning(err)
				return nil
			}
			fileNumber++
			files = append(files, info)
			fmt.Printf("%v files found...\r", fileNumber)
			return nil
		})
		fmt.Print("\n")
		if err != nil {
			return f, err
		}
		msg.Info("Calculating total size...")
		f.FilesNumber = fileNumber
		bar := progressbar.Default(fileNumber)
		for _, file := range files {
			f.Size += file.Size()
			bar.Add(1)
		}
	} else if finfo.IsDir() && !NoWalk {
		err = filepath.Walk(path, func(name string, info fs.FileInfo, err error) error {
			if err != nil {
				Warning(err)
				return nil
			}
			if PrintOnWalk {
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

	return f, nil
}

func BasicFile(finfo os.FileInfo, absPath string) (f File, err error) {
	f.Size = finfo.Size()
	f.Name = finfo.Name()
	f.IsDir = finfo.IsDir()
	f.ModTime = finfo.ModTime()
	f.Perms = finfo.Mode().Perm()
	f.AbsPath = absPath

	// Values which do not work on some systems.

	// Only on windows systems.
	f.CreationDate = stat.CreationDate(finfo)
	// Only on unix systems.
	f.User, f.Group, err = stat.UsrAndGroup(finfo)

	return
}

func BasicFileInfo(name string) (file *os.File, stat os.FileInfo, abspath string, err error) {
	file, err = os.Open(name)
	if err != nil {
		return
	}
	stat, err = file.Stat()
	if err != nil {
		return
	}
	abspath, err = filepath.Abs(name)
	if err != nil {
		return
	}
	return
}
