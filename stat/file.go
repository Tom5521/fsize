package stat

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Tom5521/fsize/filecount"
	"github.com/Tom5521/fsize/flags"
	"github.com/gookit/color"
	"github.com/labstack/gommon/bytes"
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
	_, f.info, f.linfo, f.AbsPath, err = RawInfo(path)
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

func RawInfo(name string) (file *os.File, stat, lstat os.FileInfo, abspath string, err error) {
	file, err = os.Open(name)
	if err != nil {
		return
	}
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

func (f File) String() string {
	var builder strings.Builder

	render := func(title string, content ...any) {
		fmt.Fprint(&builder, color.Green.Render(title+" "))
		fmt.Fprintln(&builder, content...)
	}

	render(po.Get("Name:"), f.Name)
	render(po.Get("Size:"), bytes.New().Format(f.Size))
	render(po.Get("Absolute path:"), f.AbsPath)
	if f.IsLink {
		render(po.Get("Physical path:"), f.PhysicalPath)
	}
	render(po.Get("Modify:"), f.ModTime.Format(time.DateTime))
	render(po.Get("Access:"), f.AccessTime.Format(time.DateTime))
	if f.SupportCreationDate {
		render(po.Get("Birth:"), f.CreationTime.Format(time.DateTime))
	}
	render(po.Get("Is directory:"), f.IsDir)
	render(po.Get("Permissions:"), fmt.Sprintf("%v/%v", int(f.Perms), f.Perms.String()))
	if f.IsDir && !flags.NoWalk {
		render(po.Get("Number of files:"), f.FilesNumber)
	}

	if f.SupportFileIDs {
		render(po.Get("UID/User:"), fmt.Sprintf("%v/%v", f.User.Uid, f.User.Username))
		render(po.Get("GID/Group:"), fmt.Sprintf("%v/%v", f.Group.Gid, f.Group.Name))
	}

	return builder.String()
}
