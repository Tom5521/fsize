//go:build windows
// +build windows

package stat

import (
	"os"
	"os/user"
)

func UsrAndGroup(info os.FileInfo) (usr *user.User, group *user.Group, err error) {
	return nil, nil, nil
}

func Usr(info os.FileInfo) (usr *user.User, err error) {
	return nil, nil
}

func Group(info os.FileInfo) (group *user.Group, err error) {
	return nil, nil
}
