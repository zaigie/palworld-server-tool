package tool

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/executor"
)

func executeCommand(command string) (*executor.Executor, string, error) {
	exec, err := executor.NewExecutor(
		viper.GetString("host"),
		viper.GetString("password"),
		viper.GetInt("timeout"), true)
	if err != nil {
		return nil, "", err
	}

	response, err := exec.Execute(command)
	return exec, response, err
}

func Info() (map[string]string, error) {
	exec, response, err := executeCommand("Info")
	if err != nil {
		return nil, err
	}
	defer exec.Close()

	re := regexp.MustCompile(`\[(v[\d\.]+)\]\s*(.+)`)
	matches := re.FindStringSubmatch(response)
	if matches == nil || len(matches) < 3 {
		return map[string]string{
			"version": "Unknown",
			"name":    "Unknown",
		}, nil
	}
	result := map[string]string{
		"version": matches[1],
		"name":    matches[2],
	}
	return result, nil
}

func ShowPlayers() ([]map[string]string, error) {
	exec, response, err := executeCommand("ShowPlayers")
	if err != nil {
		return nil, err
	}
	defer exec.Close()

	lines := strings.Split(response, "\n")
	titles := strings.Split(lines[0], ",")
	var result []map[string]string
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		playerData := make(map[string]string)
		for i, title := range titles {
			value := "<null/err>"
			if i < len(fields) {
				value = fields[i]
				if strings.Contains(value, "\u0000") {
					// Usually \u0000 is an error
					value = "<null/err>"
				}
			}
			playerData[title] = value
		}
		result = append(result, playerData)
	}
	return result, nil
}

func KickPlayer(steamID string) error {
	exec, response, err := executeCommand("KickPlayer " + steamID)
	if err != nil {
		return err
	}
	defer exec.Close()

	if response != fmt.Sprintf("Kicked: %s", steamID) {
		return errors.New(response)
	}
	return nil
}

func BanPlayer(steamID string) error {
	exec, response, err := executeCommand("BanPlayer " + steamID)
	if err != nil {
		return err
	}
	defer exec.Close()

	if response != fmt.Sprintf("Banned: %s", steamID) {
		return errors.New(response)
	}
	return nil
}

func Broadcast(message string) error {
	message = strings.ReplaceAll(message, " ", "_")
	exec, response, err := executeCommand("Broadcast " + message)
	if err != nil {
		return err
	}
	defer exec.Close()

	if response != fmt.Sprintf("Broadcasted: %s", message) {
		return errors.New(response)
	}
	return nil
}

func Shutdown(seconds string, message string) error {
	message = strings.ReplaceAll(message, " ", "_")
	exec, response, err := executeCommand(fmt.Sprintf("Shutdown %s %s", seconds, message))
	if err != nil {
		return err
	}
	defer exec.Close()

	if response != fmt.Sprintf("Shutdown: %s", message) {
		// return errors.New(response)
		return nil // HACK: Not Tested
	}
	return nil
}

func DoExit() error {
	exec, response, err := executeCommand("DoExit")
	if err != nil {
		return err
	}
	defer exec.Close()

	if response != "Exited" {
		// return errors.New(response)
		return nil // HACK: Not Tested
	}
	return nil
}
