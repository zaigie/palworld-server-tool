package executor

import (
	"errors"
	"strings"
	"time"

	"github.com/gorcon/rcon"
)

var (
	ErrPasswordEmpty = errors.New("password is empty")
)

type ExecuteCloser interface {
	Execute(command string) (string, error)
	Close() error
}

type Executor struct {
	skipErrors bool
	client     ExecuteCloser
}

func NewExecutor(address, password string, timeout int, skipErrors bool) (*Executor, error) {
	if password == "" {
		return nil, ErrPasswordEmpty
	}
	timeoutDuration := time.Duration(timeout) * time.Second
	client, err := rcon.Dial(address, password, rcon.SetDialTimeout(timeoutDuration), rcon.SetDeadline(timeoutDuration))
	if err != nil {
		return nil, err
	}
	return &Executor{client: client, skipErrors: skipErrors}, nil
}

func (e *Executor) Execute(command string) (string, error) {
	response, err := e.client.Execute(command)
	response = strings.TrimSpace(response)
	if err != nil && e.skipErrors && response != "" {
		return response, nil
	}
	return response, err
}

func (e *Executor) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}
