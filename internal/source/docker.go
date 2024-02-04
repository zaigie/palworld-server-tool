package source

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/zaigie/palworld-server-tool/internal/logger"
)

func getDockerClient() (*client.Client, error) {
	dockerAPIVersion := os.Getenv("DOCKER_API_VERSION")
	if dockerAPIVersion == "" {
		return client.NewClientWithOpts(client.FromEnv)
	} else {
		return client.NewClientWithOpts(client.FromEnv, client.WithVersion(dockerAPIVersion))
	}
}

func CopyFromContainer(containerID, remotePath string) (string, error) {
	tmpFile, err := os.CreateTemp("", "docker-file-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	execConfig := types.ExecConfig{
		Cmd:          []string{"find", remotePath, "-name", "Level.sav"},
		AttachStdout: true,
		AttachStderr: true,
	}

	ctx := context.Background()
	// dockerSocket := "unix:///app/run/docker.sock"
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(dockerSocket))
	cli, err := getDockerClient()
	if err != nil {
		return "", err
	}
	defer cli.Close()

	resp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", err
	}

	response, err := cli.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer response.Close()

	out, err := io.ReadAll(response.Reader)
	if err != nil {
		return "", err
	}

	foundFilePath := strings.TrimSpace(string(out))
	if strings.HasPrefix(foundFilePath, "K") {
		foundFilePath = strings.TrimSpace(foundFilePath[1:])
	}

	logger.Debugf("found file: %s\n", foundFilePath)
	if foundFilePath == "" {
		return "", errors.New("file Level.sav not found in container")
	}

	reader, _, err := cli.CopyFromContainer(ctx, containerID, foundFilePath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	_, err = io.Copy(tmpFile, reader)
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
