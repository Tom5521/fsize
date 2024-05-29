//go:build unix
// +build unix

package filestat

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func GetUsrAndGroup(info os.FileInfo) (usr *user.User, group *user.Group, err error) {
	usr, err = GetUsr(info)
	if err != nil {
		return
	}
	group, err = GetGroup(info)
	return
}

func GetUsr(info os.FileInfo) (usr *user.User, err error) {
	stat := info.Sys().(*syscall.Stat_t)
	usr, err = user.LookupId(strconv.Itoa(int(stat.Uid)))
	if err != nil {
		return
	}
	return
}

func GetGroup(info os.FileInfo) (group *user.Group, err error) {
	stat := info.Sys().(*syscall.Stat_t)
	group, err = user.LookupGroupId(strconv.Itoa(int(stat.Gid)))
	if err != nil {
		return
	}
	return
}
