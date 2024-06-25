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
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic(ErrGettingStruct)
	}

	usr, err = user.LookupId(strconv.FormatUint(uint64(stat.Uid), 10))
	return
}

func Group(info os.FileInfo) (group *user.Group, err error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic(ErrGettingStruct)
	}

	group, err = user.LookupGroupId(strconv.FormatUint(uint64(stat.Gid), 10))
	return
}
