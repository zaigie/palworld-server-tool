package source

import (
	"archive/tar"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/system"
)

func getDockerClient() (*client.Client, error) {
	dockerAPIVersion := os.Getenv("DOCKER_API_VERSION")
	if dockerAPIVersion == "" {
		return client.NewClientWithOpts(client.FromEnv)
	} else {
		return client.NewClientWithOpts(client.FromEnv, client.WithVersion(dockerAPIVersion))
	}
}

func CopyFromContainer(containerID, remotePath, way string) (string, error) {
	logger.Infof("copying savDir from %s\n", remotePath)
	ctx := context.Background()
	cli, err := getDockerClient()
	if err != nil {
		return "", err
	}
	defer cli.Close()

	savDir, err := getSavDir(containerID, remotePath, cli, ctx)
	if err != nil {
		return "", err
	}
	relatedSavHash := filepath.Base(savDir)

	reader, _, err := cli.CopyFromContainer(ctx, containerID, savDir)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	uuid := uuid.New().String()
	tempDir := filepath.Join(os.TempDir(), "palworldsav-docker-"+way+"-"+uuid)
	absPath, err := filepath.Abs(tempDir)
	if err != nil {
		return "", err
	}

	if err = system.CleanAndCreateDir(absPath); err != nil {
		return "", err
	}

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		path := filepath.Join(absPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return "", err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				return "", err
			}
			file.Close()
		}
	}

	levelFilePath := filepath.Join(absPath, relatedSavHash, "Level.sav")
	return levelFilePath, nil
}

func getSavDir(containerID, path string, cli *client.Client, ctx context.Context) (string, error) {
	execConfig := types.ExecConfig{
		Cmd:          []string{"find", path, "-name", "Level.sav"},
		AttachStdout: true,
		AttachStderr: true,
	}
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

	levelFilePath, err := getValidFilePath(strings.TrimSpace(string(out)), path)
	if err != nil {
		return "", err
	}

	logger.Debugf("docker find Level.sav file: %s\n", levelFilePath)
	if !strings.HasSuffix(levelFilePath, "Level.sav") {
		return "", errors.New("file Level.sav not found")
	}
	return filepath.Dir(levelFilePath), nil
}

func getValidFilePath(output, expectedStart string) (string, error) {
	startIndex := strings.Index(output, expectedStart)
	if startIndex == -1 {
		return "", errors.New("expected path not found in the output")
	}
	validPath := output[startIndex:]
	return validPath, nil
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
