package main

import (
	"os"
	"path/filepath"

	"github.com/Tom5521/fsize/checkos"
	"github.com/Tom5521/fsize/filecount"
	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/stat"
)

func Read(path string) (f stat.File, err error) {
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

	if flags.Progress && finfo.IsDir() && !flags.NoWalk {
		err = filecount.Progress(&f.FilesNumber, &f.Size, path)
	} else if finfo.IsDir() && !flags.NoWalk {
		err = filecount.Print(&f.FilesNumber, &f.Size, path)
	}

	return
}

func BasicFile(finfo os.FileInfo, absPath string) (f stat.File, err error) {
	f.Size = finfo.Size()
	f.Name = finfo.Name()
	f.IsDir = finfo.IsDir()
	f.ModTime = finfo.ModTime()
	f.Perms = finfo.Mode().Perm()
	f.AbsPath = absPath

	// Values which do not work on some systems.

	// Only on windows systems.
	if checkos.Windows {
		f.CreationDate = stat.CreationDate(finfo)
	}
	// Only on unix systems.
	if checkos.Unix {
		f.User, f.Group, err = stat.UsrAndGroup(finfo)
	}

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
