package database

import (
	"time"

	"github.com/zaigie/palworld-server-tool/internal/logger"
	"go.etcd.io/bbolt"
)

var db *bbolt.DB

func createBuckets(db *bbolt.DB) {
	// players
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("players"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}
	// guilds
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}
}

func InitDB() *bbolt.DB {
	db, err := bbolt.Open("pst.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Panic(err)
	}
	createBuckets(db)
	return db
}

func GetDB() *bbolt.DB {
	if db == nil {
		db = InitDB()
	}
	return db
}
