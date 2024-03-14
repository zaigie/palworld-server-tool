//go:build !windows

package system

import (
	"os"
	"time"
)

func GetFileCreateTime(info os.FileInfo) time.Time {
	return info.ModTime()
}

func GetEntryCreateTime(info os.DirEntry) time.Time {
	unixFileInfo, err := info.Info()
	if err != nil {
		return time.Time{}
	}
	return unixFileInfo.ModTime()
}
