package source

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CopyFromContainer(containerID, remotePath string) (string, error) {
	tmpFile, err := os.CreateTemp("", "docker-file-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	var stdout bytes.Buffer

	finderCmd := fmt.Sprintf("docker exec %s find %s -name Level.sav", containerID, remotePath)

	cmd := exec.Command("sh", "-c", finderCmd)
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	foundFilePath := strings.TrimSpace(stdout.String())
	if foundFilePath == "" {
		return "", errors.New("file Level.sav not found in container")
	}

	dockerCmd := fmt.Sprintf("docker cp %s:%s %s", containerID, foundFilePath, tmpFile.Name())

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
