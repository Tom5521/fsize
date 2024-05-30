//go:build unix
// +build unix

package stat

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func UsrAndGroup(info os.FileInfo) (usr *user.User, group *user.Group, err error) {
	usr, err = Usr(info)
	if err != nil {
		return
	}
	group, err = Group(info)
	return
}

func Usr(info os.FileInfo) (usr *user.User, err error) {
	stat := info.Sys().(*syscall.Stat_t)
	usr, err = user.LookupId(strconv.Itoa(int(stat.Uid)))
	return
}

func Group(info os.FileInfo) (group *user.Group, err error) {
	stat := info.Sys().(*syscall.Stat_t)
	group, err = user.LookupGroupId(strconv.Itoa(int(stat.Gid)))
	return
}
