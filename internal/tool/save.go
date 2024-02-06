package tool

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/auth"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/source"
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
		// http(s)://url
		logger.Infof("downloading Level.sav from %s\n", file)
		tmpFile, err = source.DownloadFromHttp(file)
		if err != nil {
			return errors.New("error downloading file: " + err.Error())
		}
		logger.Info("Level.sav downloaded\n")
	} else if strings.HasPrefix(file, "k8s://") {
		// k8s://namespace/pod/container:remotePath
		logger.Infof("copy Level.sav from %s\n", file)
		namespace, podName, container, remotePath, err := source.ParseK8sAddress(file)
		if err != nil {
			return errors.New("error parsing k8s address: " + err.Error())
		}
		tmpFile, err = source.CopyFromPod(namespace, podName, container, remotePath)
		if err != nil {
			return errors.New("error copying file from pod: " + err.Error())
		}
	} else if strings.HasPrefix(file, "docker://") {
		// docker://containerID(Name):remotePath
		logger.Infof("copy Level.sav from %s\n", file)
		containerId, remotePath, err := source.ParseDockerAddress(file)
		if err != nil {
			return errors.New("error parsing docker address: " + err.Error())
		}
		tmpFile, err = source.CopyFromContainer(containerId, remotePath)
		if err != nil {
			return errors.New("error copying file from container: " + err.Error())
		}
	} else {
		// local file
		tmpFile, err = source.CopyFromLocal(file)
		if err != nil {
			return errors.New("error copying file to temporary directory: " + err.Error())
		}
	}
	defer os.Remove(tmpFile)

	baseUrl := "http://127.0.0.1"
	if viper.GetBool("web.tls") && !strings.HasSuffix(baseUrl, "/") {
		baseUrl = viper.GetString("web.public_url")
	}

	requestUrl := fmt.Sprintf("%s:%d/api/", baseUrl, viper.GetInt("web.port"))
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
