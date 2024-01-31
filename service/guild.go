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

func GetGuild(db *bbolt.DB, adminPlayerUid string) (database.Guild, error) {
	var guild database.Guild
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		v := b.Get([]byte(adminPlayerUid))
		if v == nil {
			return ErrNoRecord
		}
		if err := json.Unmarshal(v, &guild); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return database.Guild{}, err
	}
	return guild, nil
}
