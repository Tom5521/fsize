package stat

import (
	"os"
	"os/user"
)

type FileIDs struct {
	User  *user.User
	Group *user.Group

	SupportFileIDs bool
}

func NewFileIDs(info os.FileInfo) (fids FileIDs, err error) { return newFileIDs(info) }
