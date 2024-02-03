package service

import (
	"encoding/json"

	"github.com/zaigie/palworld-server-tool/internal/database"
	"go.etcd.io/bbolt"
)

func PutGuilds(db *bbolt.DB, guilds []database.Guild) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		for _, g := range guilds {
			v, err := json.Marshal(g)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(g.AdminPlayerUid), v); err != nil {
				return err
			}
		}
		return nil
	})
}

func ListGuilds(db *bbolt.DB) ([]database.Guild, error) {
	guilds := make([]database.Guild, 0)
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		return b.ForEach(func(k, v []byte) error {
			var guild database.Guild
			if err := json.Unmarshal(v, &guild); err != nil {
				return err
			}
			guilds = append(guilds, guild)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return guilds, nil
}

func GetGuild(db *bbolt.DB, playerUID string) (database.Guild, error) {
	var guild database.Guild
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))

		// 遍历bucket中的所有guild
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var g database.Guild
			if err := json.Unmarshal(v, &g); err != nil {
				return err
			}

			// 检查当前guild的players是否包含指定的player_uid
			for _, player := range g.Players {
				if player.PlayerUid == playerUID {
					guild = g
					return nil
				}
			}
		}
		return ErrNoRecord
	})
	if err != nil {
		return database.Guild{}, err
	}
	return guild, nil
}
