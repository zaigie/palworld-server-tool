package tool

import (
	"encoding/json"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"strconv"
)

type GameApiRest struct {
}

type ResponseInfo struct {
	Version     string `json:"version"`
	ServerName  string `json:"servername"`
	Description string `json:"description"`
}

func (g *GameApiRest) Info() (map[string]string, error) {
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

func (g *GameApiRest) ShowPlayers() ([]database.PlayerRcon, error) {
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
		id, err := strconv.Atoi(player.PlayerId)
		if err != nil {
			continue
		}
		// 临时处理游戏的BUG
		if id < 0 {
			player.PlayerId = strconv.FormatUint(uint64(uint32(id)), 10)
		}
		playerRcon := database.PlayerRcon{
			PlayerUid: player.PlayerId,
			SteamId:   player.UserId,
			Nickname:  player.Name,
		}
		playersRcon = append(playersRcon, playerRcon)
	}
	return playersRcon, nil
}

type RequestUserID struct {
	UserID string `json:"userid"`
}

func (g *GameApiRest) KickPlayer(steamID string) error {
	b, err := json.Marshal(RequestUserID{
		UserID: steamID,
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

func (g *GameApiRest) BanPlayer(steamID string) error {
	b, err := json.Marshal(RequestUserID{
		UserID: steamID,
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

func (g *GameApiRest) UnBanPlayer(steamID string) error {
	b, err := json.Marshal(RequestUserID{
		UserID: steamID,
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

func (g *GameApiRest) Broadcast(message string) error {
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

func (g *GameApiRest) Shutdown(seconds int, message string) error {
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

func (g *GameApiRest) DoExit() error {
	_, err := callApi("POST", "/v1/api/stop", nil)
	if err != nil {
		return err
	}
	return nil
}
