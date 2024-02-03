package source

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFromLocal(srcFileName string) (string, error) {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	tempDir := os.TempDir()

	dstFileName := filepath.Join(tempDir, filepath.Base(srcFileName))
	dstFile, err := os.Create(dstFileName)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", err
	}

	return dstFileName, nil
}
