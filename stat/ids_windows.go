//go:build windows
// +build windows

package stat

import (
	"os"
	"os/user"
)

func UsrAndGroup(info os.FileInfo) (usr *user.User, group *user.Group, err error) { return }

func Usr(info os.FileInfo) (usr *user.User, err error) { return }

func Group(info os.FileInfo) (group *user.Group, err error) { return }
