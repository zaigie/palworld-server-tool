package service

import (
	"encoding/json"
	"errors"
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
				var existingPlayer database.Player
				if err := json.Unmarshal(existingPlayerData, &existingPlayer); err != nil {
					return err
				}

				/// 数据合并逻辑
				if existingPlayer.Level > p.Level || (existingPlayer.Level == p.Level && existingPlayer.Exp > p.Exp) {
					if len(p.Pals) > len(existingPlayer.Pals) {
						existingPlayer.Pals = p.Pals
					}
					p = existingPlayer
				} else if len(p.Pals) < len(existingPlayer.Pals) {
					p.Pals = existingPlayer.Pals
				}

				// Rcon data already has this player
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

func isUidMatch(uid1, uid2 string) bool {
	return strings.Contains(uid1, uid2) || strings.Contains(uid2, uid1)
}

func PutPlayersRcon(db *bbolt.DB, players []database.PlayerRcon) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		for _, p := range players {

			// rcon uid may not equal to player uid but may contain
			var matchedPlayerUid string
			err := db.View(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte("players"))
				return b.ForEach(func(k, v []byte) error {
					if isUidMatch(string(k), p.PlayerUid) {
						matchedPlayerUid = string(k)
						return nil
					}
					return nil
				})
			})
			if err != nil {
				return err
			}

			if matchedPlayerUid != "" {
				p.PlayerUid = matchedPlayerUid
			}

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

func AddWhitelist(db *bbolt.DB, player database.PlayerW) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取或创建白名单bucket
		b, err := tx.CreateBucketIfNotExists([]byte("whitelist"))
		if err != nil {
			return err
		}

		// 序列化玩家数据为JSON
		playerData, err := json.Marshal(player)
		if err != nil {
			return err
		}

		// 使用 findPlayerKey 检查玩家是否已经在白名单中
		key, err := findPlayerKey(b, player)
		if err != nil {
			return err
		}

		// 如果玩家已存在，更新其信息；如果不存在，创建新的键
		if key != nil {
			// 玩家已存在，更新其信息
			if err := b.Put(key, playerData); err != nil {
				return err
			}
		} else {
			// 玩家不存在，添加新玩家
			// 生成新玩家的唯一键
			newPlayerKey := []byte(player.Name + "|" + player.SteamID + "|" + player.PlayerUID)
			if err := b.Put(newPlayerKey, playerData); err != nil {
				return err
			}
		}

		return nil
	})
}

func ListWhitelist(db *bbolt.DB) ([]database.PlayerW, error) {
	var players []database.PlayerW

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("whitelist"))
		if b == nil {
			return nil // No error, just an empty list if the bucket doesn't exist.
		}

		return b.ForEach(func(k, v []byte) error {
			var player database.PlayerW
			if err := json.Unmarshal(v, &player); err != nil {
				return err
			}
			players = append(players, player)
			return nil
		})
	})

	return players, err
}

// findPlayerKey tries to find a player in the whitelist and returns the key if found.
func findPlayerKey(b *bbolt.Bucket, player database.PlayerW) ([]byte, error) {
	var keyFound []byte
	err := b.ForEach(func(k, v []byte) error {
		var existingPlayer database.PlayerW
		if err := json.Unmarshal(v, &existingPlayer); err != nil {
			return err
		}
		if matchesCriteria(existingPlayer, player) {
			keyFound = append([]byte(nil), k...) // Make a copy of the key
			return errors.New("player found")    // Use an error to break out of the iteration early.
		}
		return nil
	})

	if err != nil && err.Error() == "player found" {
		return keyFound, nil
	}

	return nil, err
}

// RemoveWhitelist removes a player from the whitelist.
func RemoveWhitelist(db *bbolt.DB, player database.PlayerW) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("whitelist"))
		if b == nil {
			return errors.New("whitelist bucket does not exist")
		}

		key, err := findPlayerKey(b, player)
		if err != nil {
			return err
		}
		if key == nil {
			return errors.New("player not found in whitelist")
		}

		return b.Delete(key)
	})
}

// matchesCriteria checks if the given player matches the criteria.
func matchesCriteria(existingPlayer, player database.PlayerW) bool {
	// 如果PlayerUID非空且匹配，认为是同一个玩家
	if player.PlayerUID != "" && existingPlayer.PlayerUID == player.PlayerUID {
		return true
	}
	// 如果Name非空且匹配，认为是同一个玩家
	if player.Name != "" && existingPlayer.Name == player.Name {
		return true
	}
	// 如果SteamID非空且匹配，认为是同一个玩家
	if player.SteamID != "" && existingPlayer.SteamID == player.SteamID {
		return true
	}
	// 如果没有任何字段匹配，返回false
	return false
}

func PutWhitelist(db *bbolt.DB, players []database.PlayerW) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// 获取或创建白名单bucket
		b, err := tx.CreateBucketIfNotExists([]byte("whitelist"))
		if err != nil {
			return err
		}

		// 清空现有的白名单
		err = b.ForEach(func(k, v []byte) error {
			return b.Delete(k)
		})
		if err != nil {
			return err
		}

		// 遍历并添加新的玩家数据到白名单
		for _, player := range players {
			playerData, err := json.Marshal(player)
			if err != nil {
				return err
			}
			identifier := player.PlayerUID
			if identifier == "" {
				if identifier = player.SteamID; identifier == "" {
					continue
				}
			}
			if err := b.Put([]byte(identifier), playerData); err != nil {
				return err
			}
		}

		return nil
	})
}
