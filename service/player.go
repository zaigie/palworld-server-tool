package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/zaigie/palworld-server-tool/internal/database"
	"go.etcd.io/bbolt"
)

func PutPlayers(db *bbolt.DB, players []database.Player) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		for _, p := range players {
			existingPlayerData := b.Get([]byte(p.PlayerUid))
			if existingPlayerData != nil {
				// Rcon data already has this player
				var existingPlayer database.Player
				if err := json.Unmarshal(existingPlayerData, &existingPlayer); err != nil {
					return err
				}
				if existingPlayer.SteamId != "" {
					p.SteamId = existingPlayer.SteamId
				}
				emptyTime := time.Time{}
				if existingPlayer.LastOnline != emptyTime {
					p.LastOnline = existingPlayer.LastOnline
				}
			}
			v, err := json.Marshal(p)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(p.PlayerUid), v); err != nil {
				return err
			}
		}
		return nil
	})
}

func PutPlayersRcon(db *bbolt.DB, players []database.PlayerRcon) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		for _, p := range players {
			existingPlayerData := b.Get([]byte(p.PlayerUid))
			var player database.Player
			if existingPlayerData == nil {
				// player online but not in database
				player.PlayerUid = p.PlayerUid
				player.SteamId = p.SteamId
				player.Nickname = p.Nickname
				player.LastOnline = time.Now()

				v, err := json.Marshal(player)
				if err != nil {
					return err
				}
				if err := b.Put([]byte(p.PlayerUid), v); err != nil {
					return err
				}
				continue
			}

			if err := json.Unmarshal(existingPlayerData, &player); err != nil {
				return err
			}

			if player.SteamId == "" || strings.Contains(player.SteamId, "000000") {
				player.SteamId = p.SteamId
			}
			player.LastOnline = time.Now()

			v, err := json.Marshal(player)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(p.PlayerUid), v); err != nil {
				return err
			}
		}
		return nil
	})
}

func ListPlayers(db *bbolt.DB) ([]database.TersePlayer, error) {
	players := make([]database.TersePlayer, 0)
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		return b.ForEach(func(k, v []byte) error {
			if strings.Contains(string(k), "000000") {
				return nil
			}
			var player database.TersePlayer
			if err := json.Unmarshal(v, &player); err != nil {
				return err
			}
			players = append(players, player)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return players, nil
}

func GetPlayer(db *bbolt.DB, playerUid string) (database.Player, error) {
	var player database.Player
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		v := b.Get([]byte(playerUid))
		if v == nil {
			return ErrNoRecord
		}
		if err := json.Unmarshal(v, &player); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return database.Player{}, err
	}
	return player, nil
}
