//go:build windows
// +build windows

package stat

import "os"

func NewFileIDs(os.FileInfo) (fids FileIDs, err error) {
	return
}
