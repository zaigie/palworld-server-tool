package source

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/system"
)

func CopyFromLocal(src, way string) (string, error) {

	isDir, err := system.CheckIsDir(src)
	if err != nil {
		logger.Errorf("error checking if %s is a directory: %v\n", src, err)
	}

	// 获得Level.sav路径
	var levelPath string
	if isDir {
		levelPath, err = system.GetLevelSavFilePath(src)
		if err != nil {
			return "", errors.New("error finding Level.sav: \n" + err.Error())
		}
	} else {
		if filepath.Base(src) == "Level.sav" {
			levelPath = src
		} else {
			return "", errors.New("specified file is not Level.sav and source is not a directory")
		}
	}
	savDir := filepath.Dir(levelPath)

	// 创建临时目录
	randId := uuid.New().String()
	tempDir := filepath.Join(os.TempDir(), "palworldsav-"+way+"-"+randId)
	if err = os.MkdirAll(tempDir, fs.ModePerm); err != nil {
		return "", err
	}

	// 拷贝文件
	files, err := filepath.Glob(filepath.Join(savDir, "*.sav"))
	if err != nil {
		return "", err
	}
	for _, file := range files {
		dist := filepath.Join(tempDir, filepath.Base(file))
		if err = system.CopyFile(file, dist); err != nil {
			return "", err
		}
	}
	playerDir := filepath.Join(savDir, "Players")
	distPlayerDir := filepath.Join(tempDir, "Players")
	if err = system.CopyDir(playerDir, distPlayerDir); err != nil {
		return "", err
	}

	distLevelPath := filepath.Join(tempDir, "Level.sav")
	return distLevelPath, nil
}
