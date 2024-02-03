package source

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CopyFromContainer(containerID, containerPath string) (string, error) {
	tmpFile, err := os.CreateTemp("", "docker-file-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	dockerCmd := fmt.Sprintf("docker cp %s:%s %s", containerID, containerPath, tmpFile.Name())

	err = exec.Command("sh", "-c", dockerCmd).Run()
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func ParseDockerAddress(address string) (containerID, filePath string, err error) {
	address = strings.TrimPrefix(address, "docker://")

	parts := strings.SplitN(address, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid Docker address format")
	}

	containerID, filePath = parts[0], parts[1]
	return containerID, filePath, nil
}
