//go:build unix
// +build unix

package stat

import (
	"os"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
)

func NewFileIDs(info os.FileInfo) (fids FileIDs, err error) {
	fids.SupportFileIDs = runtime.GOOS != "android"
	if !fids.SupportFileIDs {
		return
	}
	fids.User, fids.Group, err = usrAndGroup(info)

	return
}

func usrAndGroup(info os.FileInfo) (usr *user.User, group *user.Group, err error) {
	usr, err = fileUsr(info)
	if err != nil {
		return
	}

	group, err = fileGroup(info)
	return
}

func fileUsr(info os.FileInfo) (*user.User, error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic(ErrGettingStruct)
	}

	return user.LookupId(formatUint(stat.Uid))
}

func fileGroup(info os.FileInfo) (*user.Group, error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic(ErrGettingStruct)
	}

	return user.LookupGroupId(formatUint(stat.Gid))
}

func formatUint(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}
