package stat

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/Tom5521/fsize/checkos"
	"github.com/Tom5521/fsize/filecount"
	"github.com/Tom5521/fsize/flags"
	"github.com/gookit/color"
	"github.com/labstack/gommon/bytes"
)

var ErrGettingStruct = errors.New("error getting the corresponding structure from fileinfo.Sys()")

type File struct {
	FileTimes

	Size int64

	Name    string
	AbsPath string
	IsDir   bool
	Perms   fs.FileMode

	Group *user.Group
	User  *user.User

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

	f.Size = f.info.Size()
	f.Name = f.info.Name()
	f.IsDir = f.info.IsDir()
	f.Perms = f.info.Mode().Perm()

	// Values which do not work on some systems.

	// Only on unix systems.
	if checkos.Unix {
		f.User, f.Group, err = UsrAndGroup(f.info)
	}

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
		fmt.Fprint(&builder, color.Green.Render(title+" "))
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

	if checkos.Unix {
		render("UID/Name:", fmt.Sprintf("%v/%v", f.User.Uid, f.User.Username))
		render("GID/Name:", fmt.Sprintf("%v/%v", f.Group.Gid, f.Group.Name))
	}

	return builder.String()
}
