package source

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/zaigie/palworld-server-tool/internal/logger"
)

func CopyFromLocal(srcFileName string) (string, error) {
	var levelFile string = srcFileName
	var err error

	dir, err := isDir(srcFileName)
	if err != nil {
		logger.Errorf("error checking if %s is a directory: %v\n", srcFileName, err)
	}

	if dir {
		levelFile, err = findLevelFile(srcFileName)
		if err != nil {
			return "", errors.New("error finding Level.sav: \n" + err.Error())
		}
	}

	srcFile, err := os.Open(levelFile)
	if err != nil {
		return "", errors.New("error opening Level.sav: \n" + err.Error())
	}
	defer srcFile.Close()

	tempDir := os.TempDir()

	dstFileName := filepath.Join(tempDir, filepath.Base(levelFile))
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

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func findLevelFile(savePath string) (string, error) {
	var foundPath string
	err := filepath.Walk(savePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "Level.sav" {
			foundPath = path
			return errors.New("file found")
		}
		return nil
	})

	if foundPath == "" {
		if err != nil && !errors.Is(err, errors.New("file found")) {
			return "", err
		}
		return "", errors.New("file Level.sav not found")
	}

	return foundPath, nil
}
