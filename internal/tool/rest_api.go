package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/database"
)

var client = &http.Client{}

func callApi(method string, api string, param []byte) ([]byte, error) {

	addr := viper.GetString("rest.address")
	user := viper.GetString("rest.username")
	pass := viper.GetString("rest.password")
	timeout := viper.GetInt("rest.timeout")

	api, err := url.JoinPath(addr, api)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(method, api, bytes.NewReader(param))
	req.SetBasicAuth(user, pass)

	client.Timeout = time.Duration(timeout) * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rest: %d %s", resp.StatusCode, b)
	}
	return b, nil
}

type ResponseInfo struct {
	Version     string `json:"version"`
	ServerName  string `json:"servername"`
	Description string `json:"description"`
}

func Info() (map[string]string, error) {
	resp, err := callApi("GET", "/v1/api/info", nil)
	if err != nil {
		return nil, err
	}
	var data ResponseInfo
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	result := map[string]string{
		"version": data.Version,
		"name":    data.ServerName,
	}
	return result, nil
}

type ResponsePlayer struct {
	Name      string  `json:"name"`
	PlayerId  string  `json:"playerId"`
	UserId    string  `json:"userId"`
	Ip        string  `json:"ip"`
	Ping      float64 `json:"ping"`
	LocationX float64 `json:"location_x"`
	LocationY float64 `json:"location_y"`
	Level     int     `json:"level"`
}

type ResponsePlayers struct {
	Players []ResponsePlayer `json:"players"`
}

func ShowPlayers() ([]database.PlayerRcon, error) {
	resp, err := callApi("GET", "/v1/api/players", nil)
	if err != nil {
		return nil, err
	}
	var data ResponsePlayers
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	playersRcon := make([]database.PlayerRcon, 0)
	for _, player := range data.Players {
		id, err := strconv.ParseUint(player.PlayerId[:8], 16, 32)
		if err != nil {
			logger.Error("Parse PlayerId fail, %s \n", err)
			continue
		}
		player.PlayerId = strconv.FormatUint(id, 10)
		playerRcon := database.PlayerRcon{
			PlayerUid: player.PlayerId,
			SteamId:   player.UserId,
			Nickname:  player.Name,
		}
		playersRcon = append(playersRcon, playerRcon)
	}
	return playersRcon, nil
}

type RequestUserId struct {
	UserId string `json:"userid"`
}

func KickPlayer(steamId string) error {
	b, err := json.Marshal(RequestUserId{
		UserId: steamId,
	})
	if err != nil {
		return err
	}
	_, err = callApi("POST", "/v1/api/kick", b)
	if err != nil {
		return err
	}
	return nil
}

func BanPlayer(steamId string) error {
	b, err := json.Marshal(RequestUserId{
		UserId: steamId,
	})
	if err != nil {
		return err
	}
	_, err = callApi("POST", "/v1/api/ban", b)
	if err != nil {
		return err
	}
	return nil
}

func UnBanPlayer(steamId string) error {
	b, err := json.Marshal(RequestUserId{
		UserId: steamId,
	})
	if err != nil {
		return err
	}
	_, err = callApi("POST", "/v1/api/unban", b)
	if err != nil {
		return err
	}
	return nil
}

type RequestBroadcast struct {
	Message string `json:"message"`
}

func Broadcast(message string) error {
	b, err := json.Marshal(RequestBroadcast{
		Message: message,
	})
	if err != nil {
		return err
	}
	_, err = callApi("POST", "/v1/api/announce", b)
	if err != nil {
		return err
	}
	return nil
}

type RequestShutdown struct {
	Waittime int    `json:"waittime"`
	Message  string `json:"message"`
}

func Shutdown(seconds int, message string) error {
	b, err := json.Marshal(RequestShutdown{
		Waittime: seconds,
		Message:  message,
	})
	if err != nil {
		return err
	}
	_, err = callApi("POST", "/v1/api/shutdown", b)
	if err != nil {
		return err
	}
	return nil
}

func DoExit() error {
	_, err := callApi("POST", "/v1/api/stop", nil)
	if err != nil {
		return err
	}
	return nil
}
