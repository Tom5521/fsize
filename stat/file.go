package stat

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Tom5521/fsize/checkos"
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

	Path    string
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
	err = f.load(path)
	if err != nil {
		return f, err
	}

	return f, err
}

func (f *File) load(path string) (err error) {
	f.Path = path
	f.info, f.linfo, f.AbsPath, err = RawInfo(path)
	if err != nil {
		return err
	}

	{
		cinfo := f.info
		if checkos.Unix {
			// wrap f.info to return absolute location in .Name
			cinfo = newAbsFileInfo(cinfo, f.AbsPath)
		}
		f.FileTimes, err = NewFileTimes(cinfo)
		if err != nil {
			return err
		}
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

	return err
}

func RawInfo(name string) (stat, lstat os.FileInfo, abspath string, err error) {
	var file *os.File
	file, err = os.Open(name)
	if err != nil {
		return stat, lstat, abspath, err
	}
	defer file.Close()
	stat, err = file.Stat()
	if err != nil {
		return stat, lstat, abspath, err
	}
	lstat, err = os.Lstat(name)
	if err != nil {
		return stat, lstat, abspath, err
	}
	abspath, err = filepath.Abs(name)
	if err != nil {
		return stat, lstat, abspath, err
	}
	return stat, lstat, abspath, err
}
