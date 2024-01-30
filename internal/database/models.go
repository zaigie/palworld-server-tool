package database

import "time"

type Pal struct {
	Level     int32    `json:"level"`
	Exp       int32    `json:"exp"`
	Hp        int32    `json:"hp"`
	MaxHp     int32    `json:"max_hp"`
	Type      string   `json:"type"`
	Gender    string   `json:"gender"`
	IsLucky   bool     `json:"is_lucky"`
	IsBoss    bool     `json:"is_boss"`
	IsTower   bool     `json:"is_tower"`
	Workspeed int32    `json:"workspeed"`
	Melee     int32    `json:"melee"`
	Ranged    int32    `json:"ranged"`
	Defense   int32    `json:"defense"`
	Rank      int32    `json:"rank"`
	Skills    []string `json:"skills"`
}

type PlayerRcon struct {
	PlayerUid  string    `json:"player_uid"`
	SteamId    string    `json:"steam_id"`
	Nickname   string    `json:"nickname"`
	LastOnline time.Time `json:"last_online"`
}

type GuildPlayer struct {
	PlayerUid string `json:"player_uid"`
	Nickname  string `json:"nickname"`
}

type TersePlayer struct {
	PlayerUid      string           `json:"player_uid"`
	Nickname       string           `json:"nickname"`
	Level          int32            `json:"level"`
	Exp            int32            `json:"exp"`
	Hp             int32            `json:"hp"`
	MaxHp          int32            `json:"max_hp"`
	ShieldHp       int32            `json:"shield_hp"`
	ShieldMaxHp    int32            `json:"shield_max_hp"`
	MaxStatusPoint int32            `json:"max_status_point"`
	StatusPoint    map[string]int32 `json:"status_point"`
	FullStomach    float64          `json:"full_stomach"`
	PlayerRcon
}

type Player struct {
	TersePlayer
	Pals []*Pal `json:"pals"`
}

type Guild struct {
	Name           string         `json:"name"`
	BaseCampLevel  int32          `json:"base_camp_level"`
	AdminPlayerUid string         `json:"admin_player_uid"`
	Players        []*GuildPlayer `json:"players"`
	BaseIds        []string       `json:"base_ids"`
}
