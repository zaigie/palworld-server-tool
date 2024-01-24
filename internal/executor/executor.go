package executor

import (
	"errors"
	"strings"
	"time"

	"github.com/gorcon/rcon"
)

var (
	ErrCommandEmpty  = errors.New("command is not set")
	ErrPasswordEmpty = errors.New("password is not set")
)

type ExecuteCloser interface {
	Execute(command string) (string, error)
	Close() error
}

type Executor struct {
	skipErrors bool
	client     ExecuteCloser
}

var timeout int = 10

func NewExecutor(address, password string, skipErrors bool) (*Executor, error) {
	var client ExecuteCloser
	var err error

	if password == "" {
		return nil, ErrPasswordEmpty
	}

	timeoutDuration := time.Duration(timeout) * time.Second

	client, err = rcon.Dial(address, password, rcon.SetDialTimeout(timeoutDuration), rcon.SetDeadline(timeoutDuration))

	if err != nil {
		return nil, err
	}

	return &Executor{client: client, skipErrors: skipErrors}, nil
}

func (e *Executor) Execute(command string) (string, error) {
	if command == "" {
		return "", ErrCommandEmpty
	}

	response, err := e.client.Execute(command)

	if response != "" {
		response = strings.TrimSpace(response)
		if err != nil && e.skipErrors {
			return response, nil
		}
	}

	return response, err
}

func (e *Executor) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}
