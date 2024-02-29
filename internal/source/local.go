package source

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/zaigie/palworld-server-tool/internal/logger"
)

func CopyFromLocal(srcFileName string) (string, error) {
	var levelDir string
	var err error
	var isDir bool

	isDir, err = checkIsDir(srcFileName)
	if err != nil {
		logger.Errorf("error checking if %s is a directory: %v\n", srcFileName, err)
	}

	if isDir {
		levelDir = srcFileName
		_, err = findLevelFile(srcFileName)
		if err != nil {
			return "", errors.New("error finding Level.sav: \n" + err.Error())
		}
	} else {
		if filepath.Base(srcFileName) == "Level.sav" {
			levelDir = filepath.Dir(srcFileName)
		} else {
			return "", errors.New("specified file is not Level.sav and source is not a directory")
		}
	}

	tempDir := filepath.Join(os.TempDir(), "palworldsav")
	absPath, err := filepath.Abs(tempDir)
	if err != nil {
		return "", err
	}

	// 复制文件前确保目标目录是干净的
	if err = cleanAndCreateDir(absPath); err != nil {
		return "", err
	}

	// 复制包含Level.sav的目录
	err = copyDir(levelDir, absPath)
	if err != nil {
		return "", err
	}
	dstFileName := filepath.Join(absPath, "Level.sav")
	return dstFileName, nil
}

func checkIsDir(path string) (bool, error) {
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

func cleanAndCreateDir(dirPath string) error {
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(dirPath); err != nil {
			return err
		}
	}
	return os.MkdirAll(dirPath, 0755)
}

func copyDir(srcDir, destDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(srcFile, destFile string) error {
	input, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}
