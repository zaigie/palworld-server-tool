//go:build windows

package system

import (
	"os"
	"syscall"
	"time"
)

func GetFileCreateTime(info os.FileInfo) time.Time {
	stat, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return time.Time{}
	}
	return time.Unix(0, stat.CreationTime.Nanoseconds())
}

func GetEntryCreateTime(info os.DirEntry) time.Time {
	winFileInfo, err := info.Info()
	if err != nil {
		return time.Time{}
	}
	stat, ok := winFileInfo.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return time.Time{}
	}
	return time.Unix(0, stat.CreationTime.Nanoseconds())
}
