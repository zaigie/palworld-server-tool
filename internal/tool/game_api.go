package tool

import (
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/database"
)

type GameApi interface {
	Info() (map[string]string, error)
	ShowPlayers() ([]database.PlayerRcon, error)
	KickPlayer(steamID string) error
	BanPlayer(steamID string) error
	UnBanPlayer(steamID string) error
	Broadcast(message string) error
	Shutdown(seconds int, message string) error
	DoExit() error
}

var instance GameApi

func GetGameApi() GameApi {
	mode := viper.GetString("task.mode")
	if instance == nil {
		if mode == "rest" {
			instance = &GameApiRest{}
		} else {
			instance = &GameApiRcon{}
		}
	}
	return instance
}
