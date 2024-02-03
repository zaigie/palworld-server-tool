package source

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zaigie/palworld-server-tool/internal/logger"
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

func CopyFromPod(namespace, podName, container, remotePath string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", errors.New("error getting in-cluster config: " + err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", errors.New("error getting clientset: " + err.Error())
	}

	tmpFile, err := os.CreateTemp("", "Level.sav")
	if err != nil {
		return "", errors.New("error creating temporary file: " + err.Error())
	}
	defer tmpFile.Close()

	if namespace == "" {
		var err error
		namespace, err = getCurrentNamespace()
		if err != nil {
			return "", errors.New("error getting current namespace: " + err.Error())
		}
	}

	if container == "" {
		// var err error
		// container, err = getFirstContainerName(clientset, namespace, podName)
		// if err != nil {
		// 	return "", errors.New("error getting first container name: " + err.Error())
		// }
		return "", ErrContainerEmpty
	}

	findCmd := []string{"sh", "-c", fmt.Sprintf("find %s -name Level.sav", remotePath)}
	filePath, err := execPodCommand(clientset, config, namespace, podName, container, findCmd)
	if err != nil {
		return "", errors.New("error executing find command: " + err.Error())
	}
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return "", errors.New("file Level.sav not found in Pod")
	}
	logger.Debugf("file path: %s\n", filePath)

	catCmd := []string{"cat", filePath}
	fileContent, err := execPodCommand(clientset, config, namespace, podName, container, catCmd)
	if err != nil {
		return "", errors.New("error executing cat command: " + err.Error())
	}

	_, err = tmpFile.WriteString(fileContent)
	if err != nil {
		return "", errors.New("error writing file content: " + err.Error())
	}

	return tmpFile.Name(), nil
}

func getCurrentNamespace() (string, error) {
	ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ns)), nil
}

// func getFirstContainerName(clientset *kubernetes.Clientset, namespace, podName string) (string, error) {
// 	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
// 	if err != nil {
// 		return "", err
// 	}

// 	if len(pod.Spec.Containers) > 0 {
// 		return pod.Spec.Containers[0].Name, nil
// 	}

// 	return "", errors.New("no containers found in the Pod")
// }

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

func ParseK8sAddress(address string) (namespace, pod, container, filePath string, err error) {
	address = strings.TrimPrefix(address, "k8s://")

	parts := strings.SplitN(address, ":", 2)
	if len(parts) != 2 {
		return "", "", "", "", errors.New("invalid address format")
	}

	pathParts := strings.Split(parts[0], "/")
	switch len(pathParts) {
	// case 1: // podname
	// 	pod = pathParts[0]
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
