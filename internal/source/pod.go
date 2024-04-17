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
	"time"

	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/system"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	ErrPodNotFound    = errors.New("pod not found")
	ErrContainerEmpty = errors.New("container empty")
	ErrAddressInvalid = errors.New("invalid save.path, eg: k8s://namespace/podName:filePath")
)

func CopyFromPod(namespace, podName, container, remotePath, way string) (string, error) {
	logger.Infof("copying savDir from %s:%s\n", container, remotePath)
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", errors.New("error getting in-cluster config: " + err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", errors.New("error getting clientset: " + err.Error())
	}

	if namespace == "" {
		var err error
		namespace, err = getCurrentNamespace()
		if err != nil {
			return "", errors.New("error getting current namespace: " + err.Error())
		}
	}

	if container == "" {
		return "", ErrContainerEmpty
	}

	findCmd := []string{"sh", "-c", fmt.Sprintf("find %s -maxdepth 4 -path '*/backup/*' -prune -o -name 'Level.sav' -print | xargs dirname", remotePath)}
	savDir, err := execPodCommand(clientset, config, namespace, podName, container, findCmd)
	if err != nil {
		return "", errors.New("error executing find command: " + err.Error())
	}
	savDir = strings.TrimSpace(savDir)
	if savDir == "" {
		return "", errors.New("directory containing Level.sav not found in Pod")
	}
	logger.Debugf("Directory path: %s\n", savDir)

	tarCmd := []string{"sh", "-c", fmt.Sprintf("cd \"%s\" && tar czf - ./*.sav ./Players/*.sav", savDir)}
	tarStream, err := execPodCommandStream(clientset, config, namespace, podName, container, tarCmd)
	if err != nil {
		return "", errors.New("error executing tar command: " + err.Error())
	}

	id := uuid.New().String()
	tempDir := filepath.Join(os.TempDir(), "palworldsav-pod-"+way+"-"+id)
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = system.UnTarGzDir(tarStream, tempDir)
	if err != nil {
		return "", err
	}

	logger.Debugf("Directory copied from pod: %s\n", tempDir)

	levelFilePath := filepath.Join(tempDir, "Level.sav")
	return levelFilePath, nil
}

func getCurrentNamespace() (string, error) {
	ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ns)), nil
}

func execPodCommand(clientset *kubernetes.Clientset, config *rest.Config, namespace, podName, container string, cmd []string) (string, error) {
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   cmd,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Container: container,
		}, scheme.ParameterCodec)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return "", err
	}

	if stderr.Len() > 0 {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}

func execPodCommandStream(clientset *kubernetes.Clientset, config *rest.Config, namespace, podName, container string, cmd []string) (io.Reader, error) {
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   cmd,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Container: container,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, err
	}

	reader, writer := io.Pipe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer writer.Close()
		err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdout: writer,
			Stderr: os.Stderr,
		})
		if err != nil {
			logger.Errorf("Stream to pod failed: %v", err)
			cancel()
		}
	}()

	return reader, nil
}

func ParseK8sAddress(address string) (namespace, pod, container, filePath string, err error) {
	address = strings.TrimPrefix(address, "k8s://")

	parts := strings.SplitN(address, ":", 2)
	if len(parts) != 2 {
		return "", "", "", "", errors.New("invalid address format")
	}

	pathParts := strings.Split(parts[0], "/")
	switch len(pathParts) {
	case 2: // podname  container
		pod, container = pathParts[0], pathParts[1]
	case 3: // namespace  podname  container
		namespace, pod, container = pathParts[0], pathParts[1], pathParts[2]
	default:
		return "", "", "", "", errors.New("invalid path format")
	}

	filePath = parts[1]
	return namespace, pod, container, filePath, nil
}
