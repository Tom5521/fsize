//go:build windows

package stat

import "os"

func newFileIDs(os.FileInfo) (fids FileIDs, err error) {
	return fids, err
}
