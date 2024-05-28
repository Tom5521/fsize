//go:build !unix
// +build !unix

package fstats

import (
	"os"
	"os/user"
)

func GetUsrAndGroup(info os.FileInfo) (usr *user.User, group *user.Group, err error) {
	return nil, nil, nil
}

func GetUsr(info os.FileInfo) (usr *user.User, err error) {
	return nil, nil
}

func GetGroup(info os.FileInfo) (group *user.Group, err error) {
	return nil, nil
}
