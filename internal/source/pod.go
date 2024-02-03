package source

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	ErrPodNotFound    = errors.New("pod not found")
	ErrAddressInvalid = errors.New("invalid save.path, eg: k8s://namespace/podName:filePath")
)

func CopyFromPod(namespace, podName, remotePath string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "Level.sav")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	cmd := []string{"sh", "-c", fmt.Sprintf("find %s -name Level.sav", remotePath)}
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: cmd,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
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

	filePath := strings.TrimSpace(stdout.String())
	if filePath == "" {
		return "", errors.New("file Level.sav not found in Pod")
	}

	kubectlCmd := fmt.Sprintf("kubectl cp %s/%s:%s %s", namespace, podName, filePath, tmpFile.Name())
	err = exec.Command("sh", "-c", kubectlCmd).Run()
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func ParseK8sAddress(address string) (namespace, pod, filePath string, err error) {
	address = strings.TrimPrefix(address, "k8s://")

	parts := strings.SplitN(address, ":", 2)
	if len(parts) != 2 {
		return "", "", "", ErrAddressInvalid
	}

	nsPodParts := strings.SplitN(parts[0], "/", 2)
	if len(nsPodParts) != 2 {
		return "", "", "", ErrAddressInvalid
	}

	namespace, pod, filePath = nsPodParts[0], nsPodParts[1], parts[1]
	return namespace, pod, filePath, nil
}
