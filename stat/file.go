package stat

import (
	"io/fs"
	"os/user"
	"time"
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
