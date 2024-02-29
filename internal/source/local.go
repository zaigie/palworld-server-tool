package source

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/system"
)

func CopyFromLocal(src string) (string, error) {
	var levelDir string
	var err error
	var isDir bool

	isDir, err = system.CheckIsDir(src)
	if err != nil {
		logger.Errorf("error checking if %s is a directory: %v\n", src, err)
	}

	if isDir {
		levelDir = src
		_, err = system.GetLevelSavFilePath(src)
		if err != nil {
			return "", errors.New("error finding Level.sav: \n" + err.Error())
		}
	} else {
		if filepath.Base(src) == "Level.sav" {
			levelDir = filepath.Dir(src)
		} else {
			return "", errors.New("specified file is not Level.sav and source is not a directory")
		}
	}

	tempDir := filepath.Join(os.TempDir(), "palworldsav")
	absPath, err := filepath.Abs(tempDir)
	if err != nil {
		return "", err
	}

	if err = system.CleanAndCreateDir(absPath); err != nil {
		return "", err
	}

	err = system.CopyDir(levelDir, absPath)
	if err != nil {
		return "", err
	}
	levelFilePath := filepath.Join(absPath, "Level.sav")
	return levelFilePath, nil
}
