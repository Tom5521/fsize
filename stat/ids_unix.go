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

func Usr(info os.FileInfo) (*user.User, error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic(ErrGettingStruct)
	}

	return user.LookupId(formatUint(stat.Uid))
}

func Group(info os.FileInfo) (*user.Group, error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic(ErrGettingStruct)
	}

	return user.LookupGroupId(formatUint(stat.Gid))
}

func formatUint(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}
