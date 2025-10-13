package stat

import (
	"os"
)

type absFileInfo struct {
	absPath string
	os.FileInfo
}

func newAbsFileInfo(info os.FileInfo, absPath string) absFileInfo {
	return absFileInfo{absPath: absPath, FileInfo: info}
}

// Returns the absolute path of the file.
func (cf absFileInfo) Name() string {
	return cf.absPath
}
