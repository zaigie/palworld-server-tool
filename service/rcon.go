package service

import (
	"encoding/json"

	"github.com/zaigie/palworld-server-tool/internal/database"
	"go.etcd.io/bbolt"
	"k8s.io/apimachinery/pkg/util/uuid"
)

func AddRconCommand(db *bbolt.DB, rcon database.RconCommand) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcons"))
		v, err := json.Marshal(rcon)
		if err != nil {
			return err
		}
		uuid := uuid.NewUUID()
		return b.Put([]byte(uuid), v)
	})
}

func PutRconCommand(db *bbolt.DB, uuid string, rcon database.RconCommand) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcons"))
		v, err := json.Marshal(rcon)
		if err != nil {
			return err
		}
		return b.Put([]byte(uuid), v)
	})
}

func ListRconCommands(db *bbolt.DB) ([]database.RconCommandList, error) {
	rcons := make([]database.RconCommandList, 0)
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcons"))
		return b.ForEach(func(k, v []byte) error {
			var rcon database.RconCommand
			if err := json.Unmarshal(v, &rcon); err != nil {
				return err
			}
			rcons = append(rcons, database.RconCommandList{
				UUID:        string(k),
				RconCommand: rcon,
			})
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return rcons, nil
}

func GetRconCommand(db *bbolt.DB, uuid string) (database.RconCommand, error) {
	var rcon database.RconCommand
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcons"))
		v := b.Get([]byte(uuid))
		if v == nil {
			return ErrNoRecord
		}
		return json.Unmarshal(v, &rcon)
	})
	return rcon, err
}

func RemoveRconCommand(db *bbolt.DB, uuid string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcons"))
		return b.Delete([]byte(uuid))
	})
}
