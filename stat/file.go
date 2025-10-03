package stat

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Tom5521/fsize/checkos"
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

func (f *File) progressUpdater(
	finished, terminatedProgress chan struct{},
	warns *[]error,
	mu *sync.Mutex,
) {
	bar := progressbar.New64(-1)

	set := func() {
		mu.Lock()
		bar.Describe(
			po.Get(
				"%d files, %d errors, total size: %s",
				f.FilesNumber,
				int64(len(*warns)),
				bytes.New().Format(f.Size),
			),
		)
		bar.Set64(f.FilesNumber)
		mu.Unlock()
	}
	for {
		select {
		case <-time.After(time.Second / 3):
			set()
		case <-finished:
			set()
			bar.Finish()
			fmt.Fprintf(os.Stderr, "\n")
			close(terminatedProgress)
			return
		}
	}
}

func (f *File) progressCount() {
	var warns []error

	finished := make(chan struct{})
	finishedProgress := make(chan struct{})
	var mu sync.Mutex

	go func() {
		if flags.NoProgress {
			return
		}
		select {
		case <-time.After(flags.ProgressDelay):
			f.progressUpdater(finished, finishedProgress, &warns, &mu)
		case <-finished:
			close(finishedProgress)
			return
		}
	}()

	err := processFiles(&f.FilesNumber, &f.Size, f.AbsPath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				mu.Lock()
				defer mu.Unlock()
				warns = append(warns, err)
				return nil
			}
			if flags.PrintOnWalk && !flags.NoProgress {
				echo.Info(po.Get("Reading \"%s\"...", path))
			}

			return nil
		}, &mu,
	)
	if !flags.NoProgress {
		close(finished)
		<-finishedProgress
	}

	if err != nil {
		warns = append(warns, err)
	}

	for _, err2 := range warns {
		echo.Warning(err2.Error())
	}
}

func NewFile(path string) (f File, err error) {
	err = f.Load(path)
	if err != nil {
		return f, err
	}

	if flags.Progress && f.info.IsDir() && !flags.NoWalk {
		f.progressCount()
	} else if f.info.IsDir() && !flags.NoWalk {
		f.printedCount()
	}

	return f, err
}

func (f *File) Load(path string) (err error) {
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
