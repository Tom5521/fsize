package stat

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Tom5521/fsize/filecount"
	"github.com/Tom5521/fsize/flags"
	po "github.com/leonelquinteros/gotext"
)

var (
	ErrGettingStruct = errors.New(
		po.Get("error getting the corresponding structure from fileinfo.Sys()"),
	)
	ErrNotSupportedOnWindows = errors.New(po.Get("not supported on windows"))
)

type File struct {
	FileTimes
	FileIDs
	FileLink

	Size int64

	Name    string
	AbsPath string
	IsDir   bool
	Perms   fs.FileMode

	// IsDir vars

	FilesNumber int64

	// Raw information.

	info  os.FileInfo
	linfo os.FileInfo
}

func NewFile(path string) (f File, err error) {
	err = f.Load(path)
	if err != nil {
		return
	}

	if flags.Progress && f.info.IsDir() && !flags.NoWalk {
		err = filecount.Progress(&f.FilesNumber, &f.Size, path)
	} else if f.info.IsDir() && !flags.NoWalk {
		err = filecount.Print(&f.FilesNumber, &f.Size, path)
	}

	return
}

func (f *File) Load(path string) (err error) {
	f.info, f.linfo, f.AbsPath, err = RawInfo(path)
	if err != nil {
		return
	}

	f.FileTimes, err = NewFileTimes(f.info)
	if err != nil {
		return
	}
	f.FileIDs, err = NewFileIDs(f.info)
	if err != nil {
		return err
	}

	f.FileLink, err = NewFileLink(f.linfo, f.AbsPath)
	if err != nil {
		return err
	}

	f.Size = f.info.Size()
	f.Name = f.info.Name()
	f.IsDir = f.info.IsDir()
	f.Perms = f.info.Mode().Perm()

	return
}

func RawInfo(name string) (stat, lstat os.FileInfo, abspath string, err error) {
	var file *os.File
	file, err = os.Open(name)
	if err != nil {
		return
	}
	defer file.Close()
	stat, err = file.Stat()
	if err != nil {
		return
	}
	lstat, err = os.Lstat(name)
	if err != nil {
		return
	}
	abspath, err = filepath.Abs(name)
	if err != nil {
		return
	}
	return
}
