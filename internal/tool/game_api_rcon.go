package tool

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"regexp"
	"strings"
)

type GameApiRcon struct {
}

func (g *GameApiRcon) Info() (map[string]string, error) {
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

func (g *GameApiRcon) ShowPlayers() ([]database.PlayerRcon, error) {
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

func (g *GameApiRcon) KickPlayer(steamID string) error {
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

func (g *GameApiRcon) BanPlayer(steamID string) error {
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

func (g *GameApiRcon) UnBanPlayer(steamID string) error {
	exec, response, err := executeCommand("UnBanPlayer " + steamID)
	if err != nil {
		return err
	}
	defer exec.Close()

	if response != fmt.Sprintf("Unbanned: %s", steamID) {
		return errors.New(response)
	}
	return nil
}

func (g *GameApiRcon) Broadcast(message string) error {
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

func (g *GameApiRcon) Shutdown(seconds int, message string) error {
	message = strings.ReplaceAll(message, " ", "_")
	exec, response, err := executeCommand(fmt.Sprintf("Shutdown %d %s", seconds, message))
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

func (g *GameApiRcon) DoExit() error {
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
