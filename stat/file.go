package stat

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Tom5521/fsize/echo"
	"github.com/Tom5521/fsize/flags"
	"github.com/labstack/gommon/bytes"
	po "github.com/leonelquinteros/gotext"
	"github.com/schollz/progressbar/v3"
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

func (f *File) printedCount() {
	processFiles(
		&f.FilesNumber,
		&f.Size,
		f.AbsPath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				echo.Warning(err.Error())
				return nil
			}

			if flags.PrintOnWalk && !flags.NoProgress {
				echo.Info(po.Get("Reading \"%s\"...", path))
			}

			return nil
		},
	)
}

func (f *File) progressCount() {
	var warns []error
	var bar *progressbar.ProgressBar

	if !flags.NoProgress {
		bar = progressbar.Default(-1)
	}
	updateBar := func() {
		if flags.NoProgress {
			return
		}
		bar.Describe(
			po.Get(
				"%v files read, %v errors, total size: %s",
				f.FilesNumber,
				int64(len(warns)),
				bytes.New().Format(f.Size),
			),
		)
	}

	err := processFiles(&f.FilesNumber, &f.Size, f.AbsPath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				warns = append(warns, err)
				return nil
			}

			updateBar()
			if flags.PrintOnWalk && !flags.NoProgress {
				echo.Info(po.Get("Reading \"%s\"...", path))
			}

			return nil
		},
	)
	if err != nil {
		warns = append(warns, err)
	}

	if !flags.NoProgress {
		err = bar.Finish()
		if err != nil {
			warns = append(warns, err)
		}
	}

	for _, err2 := range warns {
		echo.Warning(err2.Error())
	}
}

func NewFile(path string) (f File, err error) {
	err = f.Load(path)
	if err != nil {
		return
	}

	if flags.Progress && f.info.IsDir() && !flags.NoWalk {
		f.progressCount()
	} else if f.info.IsDir() && !flags.NoWalk {
		f.printedCount()
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
