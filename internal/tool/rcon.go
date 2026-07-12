package tool

import (
	"encoding/base64"

	"github.com/zaigie/palworld-server-tool/internal/config"
	"github.com/zaigie/palworld-server-tool/internal/executor"
	"github.com/zaigie/palworld-server-tool/internal/logger"
)

func executeCommand(command string) (*executor.Executor, string, error) {
	return executeCommandWithSettings(command, config.Current().Rcon)
}

func executeCommandWithSettings(command string, settings config.RconConfig) (*executor.Executor, string, error) {
	useBase64 := settings.UseBase64

	exec, err := executor.NewExecutor(
		settings.Address,
		settings.Password,
		settings.Timeout, true)
	if err != nil {
		return nil, "", err
	}

	if useBase64 {
		command = base64.StdEncoding.EncodeToString([]byte(command))
	}

	response, err := exec.Execute(command)
	if err != nil {
		return nil, "", err
	}

	if useBase64 {
		decoded, err := base64.StdEncoding.DecodeString(response)
		if err != nil {
			logger.Warnf("decode base64 '%s' error: %v\n", response, err)
			return exec, response, nil
		}
		response = string(decoded)
	}

	return exec, response, nil
}

// TestRcon runs Palworld's official read-only Info command with the supplied
// settings. It does not broadcast, mutate players, or change server state.
func TestRcon(settings config.RconConfig) (string, error) {
	exec, response, err := executeCommandWithSettings("Info", settings)
	if err != nil {
		return "", err
	}
	defer exec.Close()
	return response, nil
}

func CustomCommand(command string) (string, error) {
	exec, response, err := executeCommand(command)
	if err != nil {
		return "", err
	}
	defer exec.Close()

	return response, nil
}
