package source

import (
	"archive/tar"
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

func getValidFilePath(output, expectedStart string) (string, error) {
	startIndex := strings.Index(output, expectedStart)
	if startIndex == -1 {
		return "", errors.New("expected path not found in the output")
	}
	validPath := output[startIndex:]
	return validPath, nil
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

	foundFilePath, err := getValidFilePath(strings.TrimSpace(string(out)), remotePath)
	if err != nil {
		return "", err
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

	// reader include tar header
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return "", errors.New("file not found in tar archive")
		case err != nil:
			return "", err
		case header == nil:
			continue
		}

		if header.Typeflag == tar.TypeReg && header.Name == "Level.sav" {
			logger.Debugf("got file: %s\n", header.Name)
			_, err = io.Copy(tmpFile, tarReader)
			if err != nil {
				return "", err
			}
			return tmpFile.Name(), nil
		}
	}
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
