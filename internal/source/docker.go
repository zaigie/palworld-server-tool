package source

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/system"

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

func CopyFromContainer(containerID, remotePath, way string) (string, error) {
	logger.Infof("copying savDir from %s\n", remotePath)

	cli, err := getDockerClient()
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// 取得Level.sav所在目录
	findCmd := []string{"sh", "-c", fmt.Sprintf("find %s -maxdepth 4 -path '*/backup/*' -prune -o -name 'Level.sav' -print | xargs dirname", remotePath)}
	savDir, err := execCommand(containerID, findCmd, cli)
	if err != nil {
		return "", err
	}
	savDir = strings.TrimSpace(savDir)
	if savDir == "" {
		return "", errors.New("directory containing Level.sav not found in container")
	}

	// 压缩
	tarCmd := []string{"sh", "-c", fmt.Sprintf("cd \"%s\" && tar czf - ./*.sav ./Players/*.sav", savDir)}
	tarReader, err := execCommandStream(containerID, tarCmd, cli)
	if err != nil {
		return "", err
	}

	// 创建临时目录
	id := uuid.New().String()
	tempDir := filepath.Join(os.TempDir(), "palworldsav-docker-"+way+"-"+id)
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// 解压文件
	err = system.UnTarGzDir(tarReader, tempDir)
	if err != nil {
		return "", err
	}

	levelFilePath := filepath.Join(tempDir, "Level.sav")
	return levelFilePath, nil
}

func execCommandStream(containerID string, command []string, cli *client.Client) (io.Reader, error) {
	ctx := context.Background()
	execConfig := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	}
	ir, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return nil, err
	}
	hr, err := cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}

	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		defer hr.Close()
		_, err = stdcopy.StdCopy(writer, os.Stderr, hr.Reader)
		if err != nil {
			logger.Errorf("Stream to docker failed: %v", err)
		}
	}()

	return reader, nil
}

func execCommand(containerID string, command []string, cli *client.Client) (string, error) {
	ctx := context.Background()
	execConfig := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	}
	ir, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", err
	}
	hr, err := cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer hr.Close()

	var outBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuf, os.Stderr, hr.Reader)
	if err != nil {
		return "", err
	}

	return outBuf.String(), nil
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
