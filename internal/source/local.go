package source

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/system"
)

func CopyFromLocal(src, way string) (string, error) {
	var savDir string
	var err error
	var isDir bool

	isDir, err = system.CheckIsDir(src)
	if err != nil {
		logger.Errorf("error checking if %s is a directory: %v\n", src, err)
	}

	if isDir {
		savDir, err = system.GetSavDir(src)
		if err != nil {
			return "", errors.New("error finding Level.sav: \n" + err.Error())
		}
	} else {
		if filepath.Base(src) == "Level.sav" {
			savDir = filepath.Dir(src)
		} else {
			return "", errors.New("specified file is not Level.sav and source is not a directory")
		}
	}

	uuid := uuid.New().String()
	tempDir := filepath.Join(os.TempDir(), "palworldsav-"+way+"-"+uuid)
	absPath, err := filepath.Abs(tempDir)
	if err != nil {
		return "", err
	}

	if err = system.CleanAndCreateDir(absPath); err != nil {
		return "", err
	}

	err = system.CopyDir(savDir, absPath)
	if err != nil {
		return "", err
	}

	levelFilePath := filepath.Join(absPath, "Level.sav")
	return levelFilePath, nil
}
