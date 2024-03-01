package source

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/system"
)

func DownloadFromHttp(url, way string) (string, error) {
	logger.Infof("downloading sav.zip from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	uuid := uuid.New().String()
	tempPath := filepath.Join(os.TempDir(), "palworldsav-http-"+way+"-"+uuid)
	absPath, err := filepath.Abs(tempPath)
	if err != nil {
		return "", err
	}

	if err = system.CleanAndCreateDir(absPath); err != nil {
		return "", err
	}

	tempZipFilePath := filepath.Join(absPath, "sav.zip")
	defer os.Remove(tempZipFilePath)

	zipOut, err := os.Create(tempZipFilePath)
	if err != nil {
		return "", err
	}
	defer zipOut.Close()

	_, err = io.Copy(zipOut, resp.Body)
	if err != nil {
		return "", err
	}

	err = system.UnzipDir(tempZipFilePath, absPath)
	if err != nil {
		return "", err
	}
	levelFilePath := filepath.Join(absPath, "Level.sav")
	logger.Info("sav.zip downloaded and extracted\n")
	return levelFilePath, nil
}
