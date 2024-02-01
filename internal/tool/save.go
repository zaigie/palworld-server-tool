package tool

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/auth"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/logger"
)

type Sturcture struct {
	Players []database.Player `json:"players"`
	Guilds  []database.Guild  `json:"guilds"`
}

func getSavCli() (string, error) {
	savCliPath := viper.GetString("save.decode_path")
	if _, err := os.Stat(savCliPath); err != nil {
		return "", err
	}
	return savCliPath, nil
}

func ConversionLoading(file string) error {
	var tmpFile string
	var err error

	savCli, err := getSavCli()
	if err != nil {
		return errors.New("error getting executable path: " + err.Error())
	}

	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		logger.Infof("downloading Level.sav from %s\n", file)
		tmpFile, err = downloadFile(file)
		if err != nil {
			return errors.New("error downloading file: " + err.Error())
		}
	} else {
		tmpFile, err = copyFileToTemp(file)
		if err != nil {
			return errors.New("error copying file to temporary directory: " + err.Error())
		}
	}
	defer os.Remove(tmpFile)

	requestUrl := fmt.Sprintf("http://127.0.0.1:%d/api/", viper.GetInt("web.port"))
	tokenString, err := auth.GenerateToken()
	if err != nil {
		return errors.New("error generating token: " + err.Error())
	}
	execArgs := []string{"-f", tmpFile, "--request", requestUrl, "--token", tokenString, "--clear"}
	cmd := exec.Command(savCli, execArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return errors.New("error starting command: " + err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		return errors.New("error waiting for command: " + err.Error())
	}

	return nil
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "Level.sav")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func copyFileToTemp(srcFileName string) (string, error) {
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
