package tool

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/executor"
	"github.com/zaigie/palworld-server-tool/internal/logger"
)

func executeCommand(command string) (*executor.Executor, string, error) {
	useBase64 := viper.GetBool("rcon.is_palguard")

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
	name := matches[2]
	if strings.Contains(name, "\u0000") {
		name = strings.ReplaceAll(name, "\u0000", "")
		name += "..."
	}
	result := map[string]string{
		"version": matches[1],
		"name":    name,
	}
	return result, nil
}

func ShowPlayers() ([]database.PlayerRcon, error) {
	exec, response, err := executeCommand("ShowPlayers")
	if err != nil {
		return nil, err
	}
	defer exec.Close()

	playersRcon := make([]database.PlayerRcon, 0)

	lines := strings.Split(response, "\n")[1:]
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) < 3 {
			continue
		}
		// HACK: could be changed
		nickname := fields[0]
		playerUid := fields[1]
		steamId := fields[2]
		if strings.Contains(nickname, "\u0000") {
			logger.Warnf("nickname %s contains no-ascii, cloud be not completed\n", nickname)
			nickname = strings.ReplaceAll(nickname, "\u0000", "")
		}
		if strings.Contains(playerUid, "\u0000") || strings.Contains(playerUid, "000000") {
			logger.Warnf("%s player_uid contains no-ascii case error, will be ignored\n", nickname)
			continue
		}
		if strings.Contains(steamId, "\u0000") || strings.Contains(steamId, "000000") {
			logger.Warnf("%s steam_id contains no-ascii case error, set to empty\n", nickname)
			steamId = ""
		}
		playerRcon := database.PlayerRcon{
			Nickname:  nickname,
			PlayerUid: playerUid,
			SteamId:   steamId,
		}
		playersRcon = append(playersRcon, playerRcon)
	}

	return playersRcon, nil
}

func CheckAndKickPlayers(db *bbolt.DB, players []database.PlayerRcon) error {
	whitelist, err := service.ListWhitelist(db)
	if err != nil {
		return errors.New(err.Error())
	}
	for _, player := range players {
		if !isPlayerWhitelisted(player, whitelist) {
			// 优先使用SteamId进行操作，如果没有提供，则使用PlayerUid
			identifier := player.SteamId
			if identifier == "" {
				identifier = player.PlayerUid
			}
			if err := KickPlayer(identifier); err != nil {
				logger.Warnf("Kicked %s fail, %s \n", player.Nickname, err)
			} else {
				logger.Warnf("Kicked %s successful \n", player.Nickname)
			}
		}
	}
	return nil
}

func isPlayerWhitelisted(player database.PlayerRcon, whitelist []database.PlayerW) bool {
	for _, whitelistedPlayer := range whitelist {
		if (player.PlayerUid != "" && player.PlayerUid == whitelistedPlayer.PlayerUID) ||
			(player.SteamId != "" && player.SteamId == whitelistedPlayer.SteamID) {
			return true
		}
	}
	return false
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
	isPalguard := viper.GetBool("rcon.is_palguard")
	broadcastCmd := "pgbroadcast "
	if !isPalguard {
		message = strings.ReplaceAll(message, " ", "_")
		broadcastCmd = "Broadcast "
	}
	fullCommand := broadcastCmd + message

	// 创建一个正则表达式对象 Broadacasted! or :Broadcasted: message
	re := regexp.MustCompile(`Broad(.*)casted(!|:\s.*)?`)

	exec, response, err := executeCommand(fullCommand)
	if err != nil {
		return err
	}
	defer exec.Close()

	if !re.MatchString(response) {
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
