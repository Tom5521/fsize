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
	"github.com/Tom5521/fsize/locales"
	"github.com/gookit/color"
	"github.com/labstack/gommon/bytes"
)

var (
	po               = locales.Po
	ErrGettingStruct = errors.New(
		po.Get("error getting the corresponding structure from fileinfo.Sys()"),
	)
	ErrNotSupportedOnWindows = errors.New(po.Get("not supported on windows"))
)

type File struct {
	FileTimes
	FileIDs

	Size int64

	Name    string
	AbsPath string
	IsDir   bool
	Perms   fs.FileMode

	// IsDir vars

	FilesNumber int64

	// Raw information.

	info os.FileInfo
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
	_, f.info, f.AbsPath, err = RawInfo(path)
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

	f.Size = f.info.Size()
	f.Name = f.info.Name()
	f.IsDir = f.info.IsDir()
	f.Perms = f.info.Mode().Perm()

	return
}

func RawInfo(name string) (file *os.File, stat os.FileInfo, abspath string, err error) {
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

func (f File) String() string {
	var builder strings.Builder

	render := func(title string, content ...any) {
		fmt.Fprint(&builder, color.Green.Render(po.Get(title)+" "))
		fmt.Fprintln(&builder, content...)
	}

	render("Name:", f.Name)
	render("Size:", bytes.New().Format(f.Size))
	render("Absolute Path:", f.AbsPath)
	render("Modify:", f.ModTime.Format(time.DateTime))
	render("Access:", f.AccessTime.Format(time.DateTime))
	if f.SupportCreationDate {
		render("Birth:", f.CreationTime.Format(time.DateTime))
	}
	render("Is directory:", f.IsDir)
	render("Permissions:", fmt.Sprintf("%v/%v", int(f.Perms), f.Perms.String()))
	if f.IsDir && !flags.NoWalk {
		render("Number of files:", f.FilesNumber)
	}

	if f.SupportFileIDs {
		render("UID/User:", fmt.Sprintf("%v/%v", f.User.Uid, f.User.Username))
		render("GID/Group:", fmt.Sprintf("%v/%v", f.Group.Gid, f.Group.Name))
	}

	return builder.String()
}
