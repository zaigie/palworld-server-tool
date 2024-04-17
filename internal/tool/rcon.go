package tool

import (
	"encoding/base64"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/executor"
	"github.com/zaigie/palworld-server-tool/internal/logger"
)

func executeCommand(command string) (*executor.Executor, string, error) {
	useBase64 := viper.GetBool("rcon.use_base64")

	exec, err := executor.NewExecutor(
		viper.GetString("rcon.address"),
		viper.GetString("rcon.password"),
		viper.GetInt("rcon.timeout"), true)
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

func CustomCommand(command string) (string, error) {
	exec, response, err := executeCommand(command)
	if err != nil {
		return "", err
	}
	defer exec.Close()

	return response, nil
}
