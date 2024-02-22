package database

import "time"

type Pal struct {
	Level     int32    `json:"level"`
	Exp       int64    `json:"exp"`
	Hp        int64    `json:"hp"`
	MaxHp     int64    `json:"max_hp"`
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
	Exp            int64            `json:"exp"`
	Hp             int64            `json:"hp"`
	MaxHp          int64            `json:"max_hp"`
	ShieldHp       int64            `json:"shield_hp"`
	ShieldMaxHp    int64            `json:"shield_max_hp"`
	MaxStatusPoint int32            `json:"max_status_point"`
	StatusPoint    map[string]int32 `json:"status_point"`
	FullStomach    float64          `json:"full_stomach"`
	SaveLastOnline string           `json:"save_last_online"`
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

type PlayerW struct {
	Name      string `json:"name"`
	SteamID   string `json:"steam_id"`
	PlayerUID string `json:"player_uid"`
}

type RconCommand struct {
	Command string `json:"command"`
	Remark  string `json:"remark"`
}

type RconCommandList struct {
	UUID string `json:"uuid"`
	RconCommand
}
